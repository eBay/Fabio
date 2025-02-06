package metrics

import (
	gkm "github.com/go-kit/kit/metrics"
	prommetrics "github.com/go-kit/kit/metrics/prometheus"
	promclient "github.com/prometheus/client_golang/prometheus"
)

type PromProvider struct {
	Opts    promclient.Opts
	Buckets []float64
}

func NewPromProvider(namespace, subsystem string, buckets []float64) Provider {
	namespace = clean(namespace)
	if len(subsystem) > 0 {
		subsystem = clean(subsystem)
	}
	return &PromProvider{
		Opts: promclient.Opts{
			Namespace: namespace,
			Subsystem: subsystem,
		},
		Buckets: buckets,
	}
}

func (p *PromProvider) NewCounter(name string, labels ...string) gkm.Counter {
	copts := promclient.CounterOpts(p.Opts)
	copts.Name = clean_prom(name)
	return prommetrics.NewCounterFrom(copts, labels)
}

func (p *PromProvider) NewGauge(name string, labels ...string) gkm.Gauge {
	gopts := promclient.GaugeOpts(p.Opts)
	gopts.Name = clean_prom(name)
	return prommetrics.NewGaugeFrom(gopts, labels)
}

func (p *PromProvider) NewHistogram(name string, labels ...string) gkm.Histogram {
	hopts := promclient.HistogramOpts{
		Namespace:   p.Opts.Namespace,
		Subsystem:   p.Opts.Subsystem,
		Name:        clean_prom(name),
		Help:        p.Opts.Help,
		ConstLabels: p.Opts.ConstLabels,
		Buckets:     p.Buckets,
	}
	return prommetrics.NewHistogramFrom(hopts, labels)
}
