package main

import (
	"ostopus/head/rest"
	"ostopus/head/config"
	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	logrus.Info("starting up OStopus head")
	config.InitMetrics(":7070")
	rest.MustStartRouter(":6060")
}


func MustInc(name string) {
	newCounter := prometheus.NewCounter(prometheus.CounterOpts{Name: name,})
	if err := prometheus.Register(newCounter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			// Counter already exists, no need to register
			newCounter = are.ExistingCollector.(prometheus.Counter)
		} else {
			logrus.WithError(err).Error("unable to register metric")
		}
	}
	newCounter.Inc()
}