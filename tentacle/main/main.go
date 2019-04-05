package main

import (
	"github.com/sirupsen/logrus"
	_ "ostopus/tentacle/os"
	"ostopus/tentacle/query"
	"ostopus/tentacle/rest"
)

var (
	QueryHandler *query.QueryHandler
)

func main() {
	logrus.Info("Starting up tentacle...")

	//osHandler := os.NewOSHandler()
	//queryStore := query.NewLocalQueryStore()
	rest.StartRouter(":7070")
}
