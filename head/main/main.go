package main

import (
	"github.com/sirupsen/logrus"
	"ostopus/head/rest"
)

func main() {
	logrus.Info("starting up OStopus head...")
	rest.MustStartRouter(":6060")
}
