package config

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	// Metrics
	PingCounter *prometheus.CounterVec

	metricsAddress string
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
	PingCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ostopus_pingCounter",
		Help: "How many pings were responded by tentacles",
	}, []string{"live"})
}

func AddMetricsToHandler(name string, help string, handler http.HandlerFunc) http.HandlerFunc {
	logrus.WithFields(logrus.Fields{"name": name, "help": help}).Info("adding metric to handler")
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
