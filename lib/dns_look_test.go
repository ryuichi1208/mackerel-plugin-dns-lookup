package dns

import (
	"net"
	"reflect"
	"sync"
	"testing"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

func Test_graphGen(t *testing.T) {
	type args struct {
		labelPrefix string
		device      string
	}
	tests := []struct {
		name string
		args args
		want map[string]mp.Graphs
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := graphGen(tt.args.labelPrefix, tt.args.device); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("graphGen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDNSResult_GraphDefinition(t *testing.T) {
	type fields struct {
		Prefix string
		Name   string
		min    float64
		max    float64
		avg    float64
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]mp.Graphs
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := DNSResult{
				Prefix: tt.fields.Prefix,
				Name:   tt.fields.Name,
				min:    tt.fields.min,
				max:    tt.fields.max,
				avg:    tt.fields.avg,
			}
			if got := dr.GraphDefinition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DNSResult.GraphDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDNSResult_FetchMetrics(t *testing.T) {
	type fields struct {
		Prefix string
		Name   string
		min    float64
		max    float64
		avg    float64
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := DNSResult{
				Prefix: tt.fields.Prefix,
				Name:   tt.fields.Name,
				min:    tt.fields.min,
				max:    tt.fields.max,
				avg:    tt.fields.avg,
			}
			got, err := dr.FetchMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("DNSResult.FetchMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DNSResult.FetchMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDNS_newDNS(t *testing.T) {
	type fields struct {
		r *net.Resolver
		o options
	}
	type args struct {
		o options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DNS{
				r: tt.fields.r,
				o: tt.fields.o,
			}
			d.newDNS(tt.args.o)
		})
	}
}

func TestDNS_lookup(t *testing.T) {
	type fields struct {
		r *net.Resolver
		o options
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DNS{
				r: tt.fields.r,
				o: tt.fields.o,
			}
			if err := d.lookup(); (err != nil) != tt.wantErr {
				t.Errorf("DNS.lookup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_doQuery(t *testing.T) {
	type args struct {
		d *DNS
	}
	tests := []struct {
		name    string
		args    args
		wantRes []int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := doQuery(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("doQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("doQuery() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestDNSResult_showResult(t *testing.T) {
	type fields struct {
		Prefix string
		Name   string
		min    float64
		max    float64
		avg    float64
	}
	type args struct {
		res []int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := &DNSResult{
				Prefix: tt.fields.Prefix,
				Name:   tt.fields.Name,
				min:    tt.fields.min,
				max:    tt.fields.max,
				avg:    tt.fields.avg,
			}
			dr.showResult(tt.args.res)
		})
	}
}

func Test_g(t *testing.T) {
	type args struct {
		wg *sync.WaitGroup
		d  *DNS
		c  chan<- []int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := g(tt.args.wg, tt.args.d, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("g() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		tmp []int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.tmp)
		})
	}
}

func TestDo(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Do()
		})
	}
}
