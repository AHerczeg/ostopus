package main

import (
	"OStopus/octo/rest"
	log "github.com/inconshreveable/log15"
)

func main() {
	log.Info("starting up octo")
	rest.StartRouter(":6060")
}
