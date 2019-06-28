package main

import (
	"flag"

	"github.com/AHerczeg/ostopus/head/rest"
	"github.com/sirupsen/logrus"
)

var (
	host = flag.String("host", "localhost", "the host head is serving on")
	port = flag.Int("port", 6060, "the port head is serving on")
)

func main() {
	logrus.Info("starting up OStopus head...")
	rest.StartServing(*host, *port)
}
