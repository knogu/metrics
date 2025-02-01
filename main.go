package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type gaugeCollector struct{}

var (
	myGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "enabled",
		Help: "if 1, access to youtube is enabled. if 0, not.",
	})
	ls = make([]int, 0, 10)
)

func (c gaugeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- myGauge.Desc()
}

func (c gaugeCollector) Collect(ch chan<- prometheus.Metric) {
	val := 0
	if 0 < len(ls) {
		val = 1
		ls = ls[:len(ls)-1]
	}
	ch <- prometheus.MustNewConstMetric(
		myGauge.Desc(),
		prometheus.GaugeValue,
		float64(val),
	)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ls = append(ls, 1)
}

var addr = flag.String("listen-address", ":3001", "http")

func main() {
	flag.Parse()

	var g gaugeCollector
	prometheus.MustRegister(g)

	http.HandleFunc("/enabled", handler)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
