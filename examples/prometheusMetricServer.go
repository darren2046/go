package main

import (
	. "github.com/darren2046/go"
)

func main() {
	Lg.Trace("Started")

	prom := Tools.PrometheusMetricServer("0.0.0.0:52763")
	n := prom.NewGauge("network_error", "")   // 0 for normal，1 for error
	e := prom.NewGauge("endpoint_exists", "") // 0 for not exists，1 for exists

	for {
		status, output := Cmd.GetStatusOutput("ssh 8.8.8.8 kubectl get endpoints -n namespace service-svc")
		if status != 0 {
			n.Set(1)
			Lg.Trace("Error:", output.S)
		} else {
			n.Set(0)
			if output.Has("10.244.3.4") {
				Lg.Trace("Endpoint exists")
				e.Set(1)
			} else {
				Lg.Trace("Endpoint not exists")
				e.Set(0)
			}
		}
		Time.Sleep(12)
	}
}
