package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"ostopus/head/rest"
)

var (
	address = flag.String("address", ":6060", "the address head is serving on")
)

func main() {
	logrus.Info("starting up OStopus head...")
	rest.StartServing(*address)
}
