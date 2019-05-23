package head

import (
	"flag"
	"github.com/AHerczeg/ostopus/head/rest"
	"github.com/sirupsen/logrus"
)

var (
	address = flag.String("address", ":6060", "the address head is serving on")
)

func main() {
	logrus.Info("starting up OStopus head...")
	rest.StartServing(*address)
}
