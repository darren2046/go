package golanglibs

import (
	"testing"
)

func TestPrometheusMetricServer(t *testing.T) {
	p := getPrometheusMetricServer("0.0.0.0:9301")
	c := p.NewCounter("test_counter", "this is a test Counter metrics")
	g := p.NewGauge("test_gauge", "this is a test Gauge metrics")

	go func() {
		for {
			c.Add(1)
			Time.Sleep(1)
		}
	}()

	for {
		g.Set(Float64(Random.Int(0, 10000)))
		Time.Sleep(1)
	}
}
