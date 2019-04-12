package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"ostopus/head/tentacles"
	"ostopus/shared/helpers"
	"ostopus/shared/tentacle"
)

type pingResponse struct {
	tentacle string
	response bool
}

func MustStartRouter(address string) {
	if err := StartRouter(address); err != nil {
		panic(err)
	}
}

func StartRouter(address string) error {
	logrus.Info("Starting up router")
	router := mux.NewRouter()
	setupRouter(router)
	logrus.Info("Listening and serving HTTP", "Address", address)
	if err := http.ListenAndServe(address, router); err != nil {
		logrus.WithError(err)
		return err
	}
	return nil
}

func setupRouter(router *mux.Router) {
	router.HandleFunc("/register", registerTentacle).Methods("POST")
	router.HandleFunc("/ping", pingAll).Methods("GET")
}

func registerTentacle(w http.ResponseWriter, r *http.Request) {
	var tentacle tentacle.Tentacle
	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		helpers.WriteResponse(w, http.StatusBadRequest, []byte("failed to parse tentacle"))
		return
	}
	if tentacles.Tentacles().HasTentacle(tentacle.Name) {
		helpers.WriteResponse(w, http.StatusConflict, []byte("name already in use"))
		return
	}
	tentacles.Tentacles().SaveTentacle(tentacle)
	logrus.WithFields(logrus.Fields{"Name": tentacle.Name, "Address": tentacle.Address}).Info("New tentacle registered")
	helpers.WriteResponse(w, http.StatusCreated, []byte{})
}

func pingAll(w http.ResponseWriter, _ *http.Request) {
	results := make(map[string]bool)
	allTentacles := tentacles.Tentacles().GetAllTentacles()

	// If there are no tentacles we can quit early
	if len(allTentacles) == 0 {
		helpers.WriteResponse(w, http.StatusOK, []byte("{}"))
		return
	}

	logrus.WithFields(logrus.Fields{"tentacles": len(allTentacles)}).Info("Pinging all tentacles")

	responses := make(chan pingResponse, len(allTentacles))

	for _, tentacle := range allTentacles {
		go pingTentacle(tentacle, responses)
	}


	for range allTentacles {
		r := <- responses
		results[r.tentacle] = r.response
	}

	logrus.Info("Finished pinging")


	marshaledResults, err := json.Marshal(results)
	if err != nil {
		logrus.Error(err)
		helpers.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while preparing response"))
	}
	helpers.WriteResponse(w, http.StatusOK, marshaledResults)
}

func sendQuery(url string, query string) []byte {
	marshaledCommand, err := json.Marshal(query)
	if err != nil {
		// TODO
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshaledCommand))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		// TODO
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body
}

func pingTentacle(tentacle tentacle.Tentacle, response chan<- pingResponse) {
	logrus.WithFields(logrus.Fields{"name": tentacle.Name, "address": tentacle.Address}).Info("Pinging tentacle")
	res, err := http.Get(tentacle.Address + "/ping")
	if err != nil {
		logrus.WithError(err)
		response <- pingResponse{
			tentacle: tentacle.Name,
			response: false,
		}
		return
	}

	logrus.WithFields(logrus.Fields{"name": tentacle.Name, "address": tentacle.Address, "code": res.StatusCode}).Info("Received ping response")
	response <- pingResponse{
		tentacle: tentacle.Name,
		response: res.StatusCode == http.StatusOK,
	}

	logrus.WithFields(logrus.Fields{"name": tentacle.Name, "address": tentacle.Address}).Info("Finished pinging tentacle")
}
