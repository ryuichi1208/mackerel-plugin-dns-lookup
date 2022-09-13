package dns

import (
	"context"
	"fmt"
	"net"
	"os"
	"sort"
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
	Domain   string `short:"d" long:"domain" description:"Domain Name" required:"false"`
	Type     string `short:"t" long:"type" description:"Record Type" required:"false"`
	Count    int    `short:"n" long:"count" description:"query num" required:"false" default:"3"`
	Timeout  int    `long:"timeout" description:"deadline time" required:"false"`
	Protocol string `long:"protocol" description:"Record Type" required:"false" default:"udp"`
	Thread   int    `long:"threads" description:"thread num" required:"false" default:"1"`
	Version  bool   `long:"version" description:"show version"`
	Debug    bool   `long:"debug" description:""`
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
			},
		},
	}

	return graphs
}

func (dr DNSResult) GraphDefinition() map[string]mp.Graphs {
	return graphGen("dns", dr.Name)
}

func (dr DNSResult) FetchMetrics() (map[string]float64, error) {
	m := make(map[string]float64)

	m["avg"] = dr.avg
	m["min"] = dr.min
	m["max"] = dr.max

	return m, nil
}

func (d *DNS) newDNS(o options) {
	d.o = o
	d.r = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(o.Timeout),
			}
			return d.DialContext(ctx, o.Protocol, fmt.Sprintf("%s:%d", o.Server, o.Port))
		},
	}

}

func (d *DNS) lookup() error {
	ip, err := d.r.LookupHost(context.Background(), d.o.Domain)
	if err != nil {
		return err
	}
	if d.o.Debug {
		fmt.Println(ip[0])
	}

	return nil
}

func doQuery(d *DNS) (res []int64, err error) {
	for i := 0; i < d.o.Count; i++ {
		start := time.Now()

		if err := d.lookup(); err != nil {
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

	dr.min = float64(res[0])
	dr.max = float64(res[len(res)-1])
	dr.avg = float64(sum / int64(len(res)))

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

func Run(tmp []int64) {
	dr := &DNSResult{}
	dr.showResult(tmp)

	plugin := mp.NewMackerelPlugin(dr)
	plugin.Run()
}

func Do() {
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.Version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", Name, Version)
		os.Exit(0)
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

	Run(tmp)
}
