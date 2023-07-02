package dns

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

const Name string = "mackerel-plugin-dns-lookup"
const Version string = "0.1.0"

type options struct {
	Server   string `short:"s" long:"server" description:"DNS Server" required:"false" default:"8.8.8.8"`
	Port     int    `short:"p" long:"port" description:"query num" required:"false" default:"53"`
	Domain   string `short:"d" long:"domain" description:"Domain Name" required:"true"`
	Type     string `short:"t" long:"type" description:"Record Type" required:"false"`
	Count    int    `short:"n" long:"count" description:"query num" required:"false" default:"3"`
	Timeout  int    `long:"timeout" description:"deadline time" required:"false"`
	Protocol string `long:"protocol" description:"Record Type" required:"false" default:"udp"`
	Thread   int    `long:"threads" description:"thread num" required:"false" default:"1"`
	Version  bool   `long:"version" description:"show version"`
	Debug    bool   `long:"debug" description:""`
	Verbose  bool   `long:"verbose" description:""`
}

type DNS struct {
	r *net.Resolver
	o options
}

type DNSResult struct {
	Prefix string
	Name   string
	min    float64
	max    float64
	avg    float64
	p95    float64
	p99    float64
}

func graphGen(labelPrefix, device string) map[string]mp.Graphs {
	graphs := map[string]mp.Graphs{
		device: {
			Label: labelPrefix,
			Unit:  mp.UnitFloat,
			Metrics: []mp.Metrics{
				{Name: "min", Label: "min", Diff: false},
				{Name: "max", Label: "max", Diff: false},
				{Name: "avg", Label: "avg", Diff: false},
				{Name: "p95", Label: "p95", Diff: false},
				{Name: "p99", Label: "p99", Diff: false},
			},
		},
	}

	return graphs
}

func (dr DNSResult) GraphDefinition() map[string]mp.Graphs {
	return graphGen(fmt.Sprintf("dns-%s", strings.ReplaceAll(dr.Name, ".", "-")), "responce")
}

func (dr DNSResult) FetchMetrics() (map[string]float64, error) {
	m := make(map[string]float64)

	m["avg"] = dr.avg
	m["min"] = dr.min
	m["max"] = dr.max
	m["p95"] = dr.p95
	m["p99"] = dr.p99

	return m, nil
}

func (d *DNS) newDNS(o options) {
	d.o = o
	d.r = &net.Resolver{
		// PreferGo controls whether Go's built-in DNS resolver is preferred
		PreferGo:     true,
		StrictErrors: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(o.Timeout),
			}
			return d.DialContext(ctx, o.Protocol, fmt.Sprintf("%s:%d", o.Server, o.Port))
		},
	}

}

func (d *DNS) lookup() error {
	var ip []string
	var err error
	switch d.o.Type {
	case "a", "A":
		ip, err = d.r.LookupHost(context.Background(), d.o.Domain)
	case "ptr", "PTR":
		ip, err = d.r.LookupAddr(context.Background(), d.o.Domain)
	case "txt", "TXT":
		ip, err = d.r.LookupTXT(context.Background(), d.o.Domain)
	case "cname", "CNAME":
		_, err = d.r.LookupCNAME(context.Background(), d.o.Domain)
	default:
		ip, err = d.r.LookupHost(context.Background(), d.o.Domain)
	}
	if err != nil {
		return err
	}
	if d.o.Debug && d.o.Verbose {
		for _, t := range ip {
			fmt.Printf("[DEBUG] %s\n", t)
		}
	}

	return nil
}

func doQuery(d *DNS) (res []int64, err error) {
	for i := 0; i < d.o.Count; i++ {
		start := time.Now()

		if err := d.lookup(); err != nil {
			if err != nil {
				fmt.Println(err)
			}
			return res, err
		}
		res = append(res, time.Since(start).Milliseconds())
	}

	return res, nil
}

func (dr *DNSResult) showResult(res []int64) {
	sort.Slice(res, func(i, j int) bool {
		return res[i] < res[j]
	})

	var sum int64
	for _, x := range res {
		sum += x
	}

	if len(res) > 0 {
		dr.min = float64(res[0])
		dr.max = float64(res[len(res)-1])
		dr.avg = float64(sum / int64(len(res)))
		dr.p95 = percentileN(res, 95)
		dr.p99 = percentileN(res, 99)
	}

}

func percentileN(res []int64, p int) float64 {
	i := len(res)*p/100 - 1
	if i == -1 {
		return float64(0)
	}
	return float64(res[i])
}

func g(wg *sync.WaitGroup, d *DNS, c chan<- []int64) (err error) {
	defer wg.Done()
	res, err := doQuery(d)
	if err != nil {
		return err
	}
	c <- res
	return nil
}

func intToFloat64(i []int64) sort.Float64Slice {
	f := make([]float64, len(i))
	for n := range i {
		f[n] = float64(i[n])
	}
	return f
}

func printPercentileN(numbers *sort.Float64Slice, l, n int) float64 {
	i := l*n/100 - 1
	ns := *numbers
	return ns[i]
}

func Run(tmp []int64, opts options) {
	dr := &DNSResult{
		Name: opts.Domain,
	}
	dr.showResult(tmp)

	plugin := mp.NewMackerelPlugin(dr)
	plugin.Run()
}

func Do() {
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", Name, Version)
		os.Exit(0)
	}

	if opts.Debug {
		fmt.Printf("[DEBUG] %v\n", opts)
		fmt.Printf("[DEBUG] QueryNum = %d\n", opts.Count*opts.Thread*2)
	}

	d := &DNS{}
	d.newDNS(opts)

	var wg sync.WaitGroup
	c := make(chan []int64, opts.Thread)
	for i := 0; i < opts.Thread; i++ {
		wg.Add(1)
		go g(&wg, d, c)
	}
	wg.Wait()
	close(c)

	tmp := []int64{}
	for i := range c {
		tmp = append(tmp, i...)
	}

	Run(tmp, opts)
}
