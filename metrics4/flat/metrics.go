package flat

//import (
//	"fmt"
//	"time"
//
//	"github.com/fabiolb/fabio/metrics4"
//	"github.com/fabiolb/fabio/metrics4/names"
//)
//
//type Provider struct{
//	Metrics []*metrics4.Metric
//}
//
//func (p *Provider) NewCounter(name string, labels ...string) metrics4.Counter {
//	return &Counter{Name: names.Flatten(name, labels, names.DotSeparator)}
//}
//
//func (p *Provider) NewGauge(name string, labels ...string) metrics4.Gauge {
//	return &Gauge{Name: names.Flatten(name, labels, names.DotSeparator)}
//}
//
//func (p *Provider) NewTimer(name string, labels ...string) metrics4.Timer {
//	timer := Timer{Name: names.Flatten(name, labels, names.DotSeparator)}
//
//	p.Metrics = append(p.Metrics, &timer)
//
//	return &timer
//}
//
//func (p *Provider) GetMetrics() []*metrics4.Metric {
//	return make([]*metrics4.Metric, 0)
//}
//
//func (p *Provider) Unregister(interface{}) {}
//
//type Counter struct {
//	Name string
//}
//
//func (c *Counter) Count(n int) {
//	fmt.Printf("%s:%d|c\n", c.Name, n)
//}
//
//type Gauge struct {
//	Name string
//}
//
//func (g *Gauge) Update(n int) {
//	fmt.Printf("%s:%d|g\n", g.Name, n)
//}
//
//type Timer struct {
//	Name string
//}
//
//func (t *Timer) Update(d time.Duration) {
//	fmt.Printf("%s:%d|ms\n", t.Name, d/time.Millisecond)
//}
