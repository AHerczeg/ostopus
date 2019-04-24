package main

import (
	"github.com/sirupsen/logrus"
	"ostopus/tentacle/local"
	"ostopus/tentacle/os"
	_ "ostopus/tentacle/os"
	"ostopus/tentacle/query"
	"ostopus/tentacle/rest"
)

var (
	QueryHandler *query.StdHandler
)

func main() {
	logrus.Info("Starting up tentacle...")

	local.InitSelf("test", "http://localhost:7070")

	osHandler := os.NewOSHandler()
	queryStore := query.NewLocalQueryStore()
	query.InitQueryHandler(queryStore, osHandler)

	rest.StartRouter(":7070")

}
