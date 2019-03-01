package main

import "github.com/sirupsen/logrus"

type result struct {
	Arguments string `json:"arguments"`
	Device    string `json:"device"`
	Path      string `json:"path"`
	Version   string `json:"version"`
}

func main() {
	logrus.Info("starting up tentacle")
}
