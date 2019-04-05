package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"ostopus/head/config"
	"ostopus/head/tentacles"
	"ostopus/shared/helpers"
	"ostopus/shared/tentacle"
	"sync"
)


type pingResponse struct {
	tentacle 	string
	response	bool
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
	router.HandleFunc("/pingResponse", pingAll).Methods("GET")
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
	logrus.Info("Pinging all tentacles")
	results := make(map[string]bool)
	allTentacles := tentacles.Tentacles().GetAllTentacles()

	responses := make(chan pingResponse)

	var wg sync.WaitGroup
	wg.Add(len(allTentacles))

	for _, tentacle := range allTentacles {
		go pingTentacle(tentacle, responses)
	}

	wg.Wait()

	for response := range responses {
		results[response.tentacle] = response.response
	}

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
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body
}

func pingTentacle(tentacle tentacle.Tentacle, responses chan pingResponse) {
	res, err := http.Get(tentacle.Address + "/ping")
	if err != nil {
		logrus.WithError(err)
		responses <- pingResponse{
			tentacle: tentacle.Name,
			response: false,
		}
		return
	}

	responses <- pingResponse{
		tentacle: tentacle.Name,
		response: res.StatusCode == http.StatusOK,
	}
}
