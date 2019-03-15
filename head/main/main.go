package main

import (
	"ostopus/head/rest"
	"ostopus/head/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting up OStopus head")
	config.InitMetrics(":7070")
	rest.MustStartRouter(":6060")
}
