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
)

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
	router.HandleFunc("/register", config.AddMetricsToHandler("registerTentacle", "", registerTentacle)).Methods("POST")
	router.HandleFunc("/ping", config.AddMetricsToHandler("pingAll", "Ping all tentacles", pingAll)).Methods("GET")
}

func registerTentacle(w http.ResponseWriter, r *http.Request) {
	var tentacle tentacle.Tentacle
	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		helpers.WriteResponse(w, 400, []byte("failed to parse tentacle"))
		return
	}
	if tentacles.Tentacles().HasTentacle(tentacle.Name) {
		helpers.WriteResponse(w, 409, []byte("name already in use"))
		return
	}
	tentacles.Tentacles().SaveTentacle(tentacle)
	logrus.WithFields(logrus.Fields{"Name": tentacle.Name, "Address": tentacle.Address}).Info("New tentacle registered")
	helpers.WriteResponse(w, 201, []byte{})
}

func pingAll(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Pinging all tentacles")
	results := make(map[string]bool)
	for _, tentacle := range tentacles.Tentacles().GetAllTentacles() {
		pingTentacle(tentacle, results)
	}
	marshalledResults, err := json.Marshal(results)
	if err != nil {
		// TODO
	}
	helpers.WriteResponse(w, 200, marshalledResults)
}

func sendQuery(url string, query string) []byte {
	marshalledCommand, err := json.Marshal(query)
	if err != nil {
		// TODO
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(marshalledCommand))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body
}

func pingTentacle(tentacle tentacle.Tentacle, results map[string]bool) {
	res, err := http.Get(tentacle.Address + "/ping")
	if err != nil {
		logrus.WithError(err)
	}
	if res.StatusCode == 200 {
		config.PingCounter.WithLabelValues("live").Inc()
		results[tentacle.Name] = true
	} else {
		config.PingCounter.WithLabelValues("dead").Inc()
		results[tentacle.Name] = false
	}
}
