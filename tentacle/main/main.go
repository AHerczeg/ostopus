package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"ostopus/tentacle/local"
	"ostopus/tentacle/os"
	_ "ostopus/tentacle/os"
	"ostopus/tentacle/query"
	"ostopus/tentacle/rest"
)

var (
	address		= flag.String("address", ":7070", "the address tentacle is serving on")
	name		= flag.String("selfName", "test", "the name of this specific tentacle instance")
)

func main() {
	logrus.Info("Starting up tentacle...")

	osHandler := os.NewOSHandler()
	queryStore := query.NewLocalQueryStore()
	query.InitQueryHandler(queryStore, osHandler)

	rest.StartServing(*address)

	local.InitSelf(*name, *address)

}
