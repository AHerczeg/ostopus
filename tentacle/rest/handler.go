package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/inconshreveable/log15"
	"net/http"
	"ostopus/shared/helpers"
	"ostopus/shared/tentacle"
)

func StartRouter(address string) {
	log.Info("Starting up router")
	router := mux.NewRouter()
	setupRouter(router)
	log.Info("Listening and serving HTTP", "Address", address)
	http.ListenAndServe(address, router)
}

func setupRouter(router *mux.Router) {
	router.HandleFunc("/query", receiveCommand).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/ping", ping).Methods("GET")
}

func register(writer http.ResponseWriter, r *http.Request) {

}

func receiveCommand(w http.ResponseWriter, r *http.Request) {
	var tentacle tentacle.Tentacle
	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		helpers.WriteResponse(w, 400, []byte("failed to parse command"))
		return
	}
}

func ping (w http.ResponseWriter, _ *http.Request)  {
	w.WriteHeader(http.StatusOK)
}
