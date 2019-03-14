package rest

import (
	"github.com/gorilla/mux"
	log "github.com/inconshreveable/log15"
	"net/http"
)

func StartRouter(address string) {
	log.Info("Starting up router")
	router := mux.NewRouter()
	setupRouter(router)
	log.Info("Listening and serving HTTP", "Address", address)
	http.ListenAndServe(address, router)
}

func setupRouter(router *mux.Router) {
	router.HandleFunc("/", receiveCommand).Methods("POST")
	router.HandleFunc("/register/{address}", register).Methods("GET")
}

func register(writer http.ResponseWriter, r *http.Request) {

}

func receiveCommand(w http.ResponseWriter, r *http.Request) {

}
