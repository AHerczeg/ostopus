package main

import (
	"OStopus/tentacle/os"
	"OStopus/tentacle/query"
	"fmt"
	"github.com/sirupsen/logrus"
)

type result struct {
	Arguments string `json:"arguments"`
	Device    string `json:"device"`
	Path      string `json:"path"`
	Version   string `json:"version"`
}

func main() {
	logrus.Info("starting up tentacle")
	osHandler := os.NewOSHandler()
	queryStore := query.NewQueryStore()
	queryHandler := query.NewQueryHandler(queryStore, osHandler)

	fmt.Println(queryHandler.RunSavedQuery("kernel_info"))
}
