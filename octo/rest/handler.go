package rest

import (
	"OStopus/octo/tentacles"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/inconshreveable/log15"
	"io/ioutil"
	"net/http"
)

func StartRouter(address string) {
	log.Info("Starting up router")
	router := mux.NewRouter()
	setupRouter(router)
	log.Info("Listening and serving HTTP", "Address", address)
	http.ListenAndServe(address, router)
}k


func setupRouter(router *mux.Router) {
	router.HandleFunc("/register", registerTentacle).Methods("POST")
	router.HandleFunc("/ping", pingAll).Methods("GET")
}

func registerTentacle(w http.ResponseWriter, r *http.Request)  {
	var tentacle tentacles.Tentacle
	_ = json.NewDecoder(r.Body).Decode(&tentacle)
	tentacles.Tentacles().SaveTentacle(tentacle)
	log.Info("New tentacle registered", "Name", tentacle.Name, "Address", tentacle.Address)
	json.NewEncoder(w).Encode(tentacle)
}

func sendCommand(w http.ResponseWriter, r *http.Request)  {

}

func pingAll(w http.ResponseWriter, r *http.Request)  {
	log.Info("Pinging all tentacles")
	for _, tentacle := range tentacles.Tentacles().GetAllTentacles() {
		go pingTentacle(tentacle)
	}
}

func sendToURL(url, command string) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(command)))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body
}

func pingTentacle(tentacle tentacles.Tentacle) {
	fmt.Printf("Pinging %s: %v", tentacle.Name, sendToURL(tentacle.Address, "ping"))
}