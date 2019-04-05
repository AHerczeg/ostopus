package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/inconshreveable/log15"
	"io/ioutil"
	"net/http"
	"net/url"
	"ostopus/shared/helpers"
	"ostopus/shared/tentacle"
	"ostopus/tentacle/local"
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

func register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading request body",
			http.StatusInternalServerError)
	}
	headAddress := string(body)
	if _, err = url.ParseRequestURI(headAddress); err != nil {
		http.Error(w, "error parsing head address",
			http.StatusBadRequest)
		return
	}

	self := local.GetSelf()

	marshaledSelf, err := json.Marshal(self)
	if err != nil {
		http.Error(w, "unexpected error",
			http.StatusInternalServerError)
	}

	req, err := http.NewRequest("POST", headAddress, bytes.NewBuffer(marshaledSelf))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "unexpected error while registering tentacle",
			http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	helpers.WriteResponse(w, resp.StatusCode, respBody)
}

func receiveCommand(w http.ResponseWriter, r *http.Request) {
	var tentacle tentacle.Tentacle
	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		http.Error(w, "error parsing command",
			http.StatusBadRequest)
		return
	}
	// TODO
}

func ping (w http.ResponseWriter, _ *http.Request)  {
	w.WriteHeader(http.StatusOK)
}
