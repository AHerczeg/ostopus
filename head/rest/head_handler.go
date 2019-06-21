package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/AHerczeg/ostopus/head/tentacles"
	"github.com/AHerczeg/ostopus/shared"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	_ "github.com/AHerczeg/ostopus/head/docs"
	"github.com/swaggo/http-swagger"
)

type pingResponse struct {
	tentacle string
	response bool
}

type queryRequest struct {
	Targets []string
	Query   shared.Query
}

func StartServing(address string) {
	MustStartRouter(address, setupRouter)
}

func setupRouter(router *mux.Router) {
	router.HandleFunc("/register", registerTentacle).Methods("POST")
	router.HandleFunc("/remove", removeTentacle).Methods("DELETE")

	// swagger:operation GET /ping ping
	//
	// ---
	// responses:
	//   '200':
	//     description: successful operation
	router.HandleFunc("/ping", pingAll).Methods("GET")

	router.HandleFunc("/query", relayQuery).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFound)

}

func registerTentacle(w http.ResponseWriter, r *http.Request) {
	var tentacle shared.Tentacle

	if err := json.NewDecoder(r.Body).Decode(&tentacle); err != nil {
		writeResponse(w, http.StatusBadRequest, []byte("failed to parse tentacle"))
		return
	}

	if tentacles.Tentacles().HasTentacle(tentacle.Name) {
		writeResponse(w, http.StatusConflict, []byte("name already in use"))
		return
	}

	tentacles.Tentacles().SaveTentacle(tentacle)
	logrus.WithFields(logrus.Fields{"Name": tentacle.Name, "Address": tentacle.Address}).Info("New tentacle registered")
	writeResponse(w, http.StatusCreated, []byte{})
}

func removeTentacle(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Error(err)
		writeResponse(w, http.StatusInternalServerError, []byte("unable to read request"))
		return
	}

	result := tentacles.Tentacles().RemoveTentacle(string(request))
	if result {
		writeResponse(w, http.StatusOK, []byte{})
	} else {
		writeResponse(w, http.StatusNotFound, []byte{})
	}
}

func pingAll(w http.ResponseWriter, _ *http.Request) {
	results := make(map[string]bool)
	allTentacles := tentacles.Tentacles().GetAllTentacles()

	// If there are no tentacles we can quit early
	if len(allTentacles) == 0 {
		writeResponse(w, http.StatusOK, []byte("{}"))
		return
	}

	logrus.WithFields(logrus.Fields{"tentacles": len(allTentacles)}).Info("Pinging all tentacles")

	responses := make(chan pingResponse, len(allTentacles))

	for _, tentacle := range allTentacles {
		go pingTentacle(tentacle, responses)
	}

	for range allTentacles {
		r := <-responses
		results[r.tentacle] = r.response
	}

	logrus.Info("Finished pinging")

	marshaledResults, err := json.Marshal(results)

	if err != nil {
		logrus.Error(err)
		writeResponse(w, http.StatusInternalServerError, []byte("unexpected error while preparing response"))
	}

	writeResponse(w, http.StatusOK, marshaledResults)
}

func relayQuery(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request queryRequest
	err := decoder.Decode(&request)
	if err != nil {
		logrus.Error(err)
		writeResponse(w, http.StatusInternalServerError, []byte("unexpected error while reading request"))
	}

	logrus.WithFields(logrus.Fields{"targets": request.Targets, "query": request.Query}).Info("Relaying new query")

	if !request.Query.Validate() {
		writeResponse(w, http.StatusBadRequest, []byte("malformed query or frequency"))
	}

	marshaledQuery, err := json.Marshal(request.Query)
	if err != nil {
		logrus.Error(err)
		writeResponse(w, http.StatusInternalServerError, []byte("unexpected error while parsing query"))
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
		writeResponse(w, http.StatusInternalServerError, []byte("unexpected error while parsing results"))
	}

	writeResponse(w, http.StatusOK, marshaledMap)

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
	req, err := http.NewRequest("POST", address+"/query", bytes.NewBuffer(query))
	client := GetDefaultClient()
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "error while processing response"
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
	writeResponse(w, http.StatusNotFound, []byte{})
}

func writeResponse(w http.ResponseWriter, code int, response []byte) {
	w.WriteHeader(code)
	if len(response) > 0 {
		w.Write(response)
	}
}

func GetDefaultClient() http.Client {
	return http.Client{
		Timeout: time.Second * 10,
	}
}

func MustStartRouter(address string, routerSetup func(*mux.Router)) {
	if err := startRouter(address, routerSetup); err != nil {
		panic(err)
	}
}

func startRouter(address string, routerSetup func(*mux.Router)) error {
	logrus.Info("Starting up router")

	router := mux.NewRouter()
	router.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	routerSetup(router)
	logrus.WithFields(logrus.Fields{
		"address": address,
	}).Info("Listening and serving HTTP", "Address")

	if err := http.ListenAndServe(address, router); err != nil {
		logrus.WithError(err)
		return err
	}

	return nil
}
