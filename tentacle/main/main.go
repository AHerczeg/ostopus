package main

import (
	"github.com/sirupsen/logrus"
	_ "ostopus/tentacle/os"
	"ostopus/tentacle/query"
	"ostopus/tentacle/rest"
	"ostopus/tentacle/local"
)

var (
	QueryHandler *query.QueryHandler
)

func main() {
	logrus.Info("Starting up tentacle...")

	local.InitSelf("test", "localhost:7070")

	//osHandler := os.NewOSHandler()
	//queryStore := query.NewLocalQueryStore()
	rest.StartRouter(":7070")


}
