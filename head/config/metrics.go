package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"fmt"
)

var (
	metricsAddress string
	metrics []prometheus.Collector
)

func InitMetrics(address string) {
	populateMetrics()
	http.Handle("/metrics", promhttp.Handler())
	logrus.WithField("serving metrics", logrus.Fields{"address": address})
	go func() {
		metricsAddress = address
		if err := http.ListenAndServe(address, nil); err != nil {
			logrus.WithError(err).Error("unable to serve metrics")
			address = ""
		}
	}()
}

func populateMetrics() {
}

func AddMetricsToHandler(name string, help string, handler http.HandlerFunc) http.HandlerFunc {
	fmt.Println("AAAAAAAA")
	logrus.WithField("adding metric to handler", logrus.Fields{"name":name, "help":help})
	return promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: help,
			},
			[]string{"code"},
		),
		handler,
	)
}