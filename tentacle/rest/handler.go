package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"ostopus/shared"
	"ostopus/tentacle/local"
)

func StartRouter(address string) {
	logrus.Info("Starting up router")
	router := mux.NewRouter()
	setupRouter(router)
	logrus.Info("Listening and serving HTTP", "Address", address)
	http.ListenAndServe(address, router)
}

func setupRouter(router *mux.Router) {
	router.HandleFunc("/query", receiveCommand).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/ping", ping).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(notFound)
}

func register(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Registering tentacle")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "error reading request body",
			http.StatusInternalServerError)
	}
	headAddress := string(body)
	if _, err = url.ParseRequestURI(headAddress); err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusBadRequest, []byte("error parsing head address"))
		return
	}

	logrus.WithFields(logrus.Fields{
		"address": headAddress,
	}).Info("Received head address")

	self := local.GetSelf()

	marshaledSelf, err := json.Marshal(self)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "unexpected error",
			http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", headAddress, bytes.NewBuffer(marshaledSelf))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while registering tentacle"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		local.SetHeadAddress(headAddress)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	shared.WriteResponse(w, resp.StatusCode, respBody)
}

func receiveCommand(w http.ResponseWriter, r *http.Request) {
	var tentacle shared.Tentacle
	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		logrus.Error(err)
		http.Error(w, "error parsing command",
			http.StatusBadRequest)
		return
	}

	shared.WriteResponse(w, http.StatusNotFound, []byte{})
	// TODO
}

func ping(w http.ResponseWriter, _ *http.Request) {
	logrus.Info("Receiving ping")
	shared.WriteResponse(w, http.StatusOK, []byte{})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	shared.WriteResponse(w, http.StatusNotFound, []byte{})
}

