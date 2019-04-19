package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"ostopus/head/tentacles"
	"ostopus/shared"
	"sync"
)

type pingResponse struct {
	tentacle string
	response bool
}

type queryRequest struct {
	Targets []string
	Query   shared.Query
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
	router.HandleFunc("/remove", removeTentacle).Methods("DELETE")
	router.HandleFunc("/ping", pingAll).Methods("GET")
	router.HandleFunc("/query", relayQuery).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)

}



func registerTentacle(w http.ResponseWriter, r *http.Request) {
	var tentacle shared.Tentacle

	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		shared.WriteResponse(w, http.StatusBadRequest, []byte("failed to parse tentacle"))
		return
	}

	if tentacles.Tentacles().HasTentacle(tentacle.Name) {
		shared.WriteResponse(w, http.StatusConflict, []byte("name already in use"))
		return
	}

	tentacles.Tentacles().SaveTentacle(tentacle)
	logrus.WithFields(logrus.Fields{"Name": tentacle.Name, "Address": tentacle.Address}).Info("New tentacle registered")
	shared.WriteResponse(w, http.StatusCreated, []byte{})
}

func removeTentacle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request string
	err := decoder.Decode(&request)
	if err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while reading request"))
	}

	result := tentacles.Tentacles().RemoveTentacle(request)
	if result {
		shared.WriteResponse(w, http.StatusOK, []byte{})
	} else {
		shared.WriteResponse(w, http.StatusInternalServerError, []byte{})
	}
}

func pingAll(w http.ResponseWriter, _ *http.Request) {
	results := make(map[string]bool)
	allTentacles := tentacles.Tentacles().GetAllTentacles()

	// If there are no tentacles we can quit early
	if len(allTentacles) == 0 {
		shared.WriteResponse(w, http.StatusOK, []byte("{}"))
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
		shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while preparing response"))
	}

	shared.WriteResponse(w, http.StatusOK, marshaledResults)
}

func relayQuery(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request queryRequest
	err := decoder.Decode(&request)
	if err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while reading request"))
	}

	logrus.WithFields(logrus.Fields{"targets": request.Targets, "query": request.Query}).Info("Relaying new query")

	if !request.Query.Validate(){
		shared.WriteResponse(w, http.StatusBadRequest, []byte("malformed query or frequency"))
	}

	marshaledQuery, err := json.Marshal(request.Query)
	if err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while parsing query"))
	}

	var results sync.Map
	var wg sync.WaitGroup

	for _, target := range request.Targets {
		if tentacle, ok := tentacles.Tentacles().GetTentacle(target); ok {
			wg.Add(1)
			go syncQuery(marshaledQuery, tentacle, &results, &wg)
		} else {
			results.Store(target, "unknown target")
		}
	}

	wg.Wait()

	marshaledMap, err := marshalSyncMap(results)
	if err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while parsing results"))
	}

	shared.WriteResponse(w, http.StatusOK, marshaledMap)

}

func pingTentacle(tentacle shared.Tentacle, response chan<- pingResponse) {
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

func syncQuery(query []byte, tentacle shared.Tentacle, results *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	results.Store(tentacle.Name, sendQuery(query, tentacle.Address))
}

func sendQuery(query []byte, address string) string {
	req, err := http.NewRequest("POST", address + "/query", bytes.NewBuffer(query))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return  "error while processing response"
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("failed query. Code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "error while reading response"
	}
	return string(bodyBytes)

}

func marshalSyncMap(syncMap sync.Map) ([]byte, error) {
	tmpMap := make(map[string]string)
	syncMap.Range(func(k, v interface{}) bool {
		tmpMap[k.(string)] = v.(string)
		return true
	})
	return json.Marshal(tmpMap)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	shared.WriteResponse(w, http.StatusNotFound, []byte{})
}