package main

import (
	"github.com/sirupsen/logrus"
	"ostopus/tentacle/query"
)

var (
	QueryHandler *query.QueryHandler
)

func main() {
	logrus.Info("starting up tentacle")
	//osHandler := os.NewOSHandler()
	//queryStore := query.NewLocalQueryStore()
	//http.ListenAndServe()
}
