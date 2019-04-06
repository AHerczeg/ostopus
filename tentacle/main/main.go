package main

import (
	"github.com/sirupsen/logrus"
	"ostopus/tentacle/local"
	_ "ostopus/tentacle/os"
	"ostopus/tentacle/query"
	"ostopus/tentacle/rest"
)

var (
	QueryHandler *query.QueryHandler
)

func main() {
	logrus.Info("Starting up tentacle...")

	local.InitSelf("test", "http://localhost:7070")

	//osHandler := os.NewOSHandler()
	//queryStore := query.NewLocalQueryStore()
	rest.StartRouter(":7070")

}
