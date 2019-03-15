package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	metrics []prometheus.Collector
)

func InitMetrics() {
	populateMetrics()
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
	http.Handle("/watermelon", promhttp.Handler())
	go http.ListenAndServe(":1971", nil)
}

func populateMetrics() {
	purchaseSubCreateCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "plutus_created_subs_total",
		Help: "The total number purchase subscriptions created",
	})
	metrics = append(metrics, purchaseSubCreateCounter)


}