package main

import (
	"OStopus/tentacle/os"
	"OStopus/tentacle/query"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	QueryHandler *query.QueryHandler
)

func main() {
	logrus.Info("starting up tentacle")
	osHandler := os.NewOSHandler()
	queryStore := query.NewLocalQueryStore()
	http.ListenAndServe()
}
