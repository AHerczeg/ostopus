package main

import (
	"flag"
	"github.com/AHerczeg/ostopus/tentacle/local"
	"github.com/AHerczeg/ostopus/tentacle/os"
	_ "github.com/AHerczeg/ostopus/tentacle/os"
	"github.com/AHerczeg/ostopus/tentacle/query"
	"github.com/AHerczeg/ostopus/tentacle/rest"
	"github.com/sirupsen/logrus"
)

var (
	address = flag.String("address", ":7070", "the address tentacle is serving on")
	name    = flag.String("selfName", "test", "the name of this specific tentacle instance")
)

func main() {
	logrus.Info("Starting up tentacle...")

	osHandler := os.NewOSHandler()
	queryStore := query.NewLocalQueryStore()
	query.InitQueryHandler(queryStore, osHandler)

	rest.StartServing(*address)

	local.InitSelf(*name, *address)

}
