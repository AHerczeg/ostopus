package rest

import (
	"ostopus/head/tentacles"
	"ostopus/shared/command"
	"ostopus/shared/tentacle"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"net/http"
	"io/ioutil"
	"github.com/sirupsen/logrus"
)

func MustStartRouter(address string)  {
	if err := StartRouter(address); err != nil {
		panic(err)
	}
}

func StartRouter(address string) error{
	log.Info("Starting up router")
	router := mux.NewRouter()
	setupRouter(router)
	logrus.Info("Listening and serving HTTP", "Address", address)
	if err := http.ListenAndServe(address, router); err != nil {
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
	_ = json.NewDecoder(r.Body).Decode(&tentacle)
	tentacles.Tentacles().SaveTentacle(tentacle)
	log.Info("New tentacle registered", "Name", tentacle.Name, "Address", tentacle.Address)
	json.NewEncoder(w).Encode(tentacle)
}

func sendCommand(w http.ResponseWriter, r *http.Request) {

}

func pingAll(w http.ResponseWriter, r *http.Request) {
	log.Info("Pinging all tentacles")
	for _, tentacle := range tentacles.Tentacles().GetAllTentacles() {
		go pingTentacle(tentacle)
	}
}

func sendToURL(url string, command command.Command) []byte {
	marshalledCommand, err := json.Marshal(command)
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

func pingTentacle(tentacle tentacle.Tentacle) {
	fmt.Printf("Pinging %s: %v", tentacle.Name, sendToURL(tentacle.Address, command.Command{Name: "ping"}))
}
