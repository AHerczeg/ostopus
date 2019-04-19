package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"ostopus/shared"
	"ostopus/tentacle/local"
	tentacleQuery "ostopus/tentacle/query"
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
	var query shared.Query
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		logrus.Error(err)
		shared.WriteResponse(w, http.StatusBadRequest, []byte("failed to unmarshal query"))
		return
	}

	if query.Frequency != "" {
		// TODO handle frequency
	} else {
		if !query.Validate() {
			shared.WriteResponse(w, http.StatusBadRequest, []byte("malformed query or frequency"))
			return
		}

		queryHandler := tentacleQuery.GetQueryHandler()

		result, err := queryHandler.RunCustomQuery(query.Query)
		fmt.Println(result)
		if err != nil {
			shared.WriteResponse(w, http.StatusInternalServerError, []byte(fmt.Sprintf("failed to run query. err: %v", err)))
			return
		}

		marshaledResults, err := json.Marshal(result)
		fmt.Println(string(marshaledResults))
		if err != nil {
			logrus.Error(err)
			shared.WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while parsing result"))
		}

		shared.WriteResponse(w, http.StatusOK, marshaledResults)
	}
}

func ping(w http.ResponseWriter, _ *http.Request) {
	logrus.Info("Receiving ping")
	shared.WriteResponse(w, http.StatusOK, []byte{})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	shared.WriteResponse(w, http.StatusNotFound, []byte{})
}

