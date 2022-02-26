package golanglibs

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusMetricServerStruct struct {
}

func getPrometheusMetricServer(listenAddr string, path ...string) *prometheusMetricServerStruct {
	g := gin.New()
	var p string
	if len(path) == 0 {
		p = "/metrics"
	} else {
		p = path[0]
	}
	g.GET("/"+p, gin.WrapH(promhttp.Handler()))

	go func() {
		g.Run(listenAddr)
	}()

	return &prometheusMetricServerStruct{}
}

func (m *prometheusMetricServerStruct) NewCounter(name string, help string) prometheus.Counter {
	pcounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(pcounter)
	return pcounter
}

func (m *prometheusMetricServerStruct) NewGauge(name string, help string) prometheus.Gauge {
	pgauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(pgauge)
	return pgauge
}
