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

type PrometheusCounterStruct struct {
	p prometheus.Counter
}

func (m *PrometheusCounterStruct) Add(num float64) {
	m.p.Add(num)
}

func (m *prometheusMetricServerStruct) NewCounter(name string, help string) *PrometheusCounterStruct {
	pcounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(pcounter)
	return &PrometheusCounterStruct{p: pcounter}
}

type PrometheusCounterVecStruct struct {
	p *prometheus.CounterVec
}

func (m *PrometheusCounterVecStruct) Add(num float64, label map[string]string) {
	m.p.With(label).Add(num)
}

func (m *prometheusMetricServerStruct) NewCounterWithLabel(name string, label []string, help string) *PrometheusCounterVecStruct {
	pcounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: name,
		Help: help,
	}, label)
	prometheus.MustRegister(pcounter)
	return &PrometheusCounterVecStruct{p: pcounter}
}

type PrometheusGaugeStruct struct {
	p prometheus.Gauge
}

func (m *PrometheusGaugeStruct) Set(num float64) {
	m.p.Set(num)
}

func (m *prometheusMetricServerStruct) NewGauge(name string, help string) *PrometheusGaugeStruct {
	pgauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
	prometheus.MustRegister(pgauge)
	return &PrometheusGaugeStruct{p: pgauge}
}

type PrometheusGaugeVecStruct struct {
	p *prometheus.GaugeVec
}

func (m *PrometheusGaugeVecStruct) Set(num float64, label map[string]string) {
	m.p.With(label).Set(num)
}

func (m *prometheusMetricServerStruct) NewGaugeWithLabel(name string, label []string, help string) *PrometheusGaugeVecStruct {
	pgauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	}, label)
	prometheus.MustRegister(pgauge)
	return &PrometheusGaugeVecStruct{p: pgauge}
}
