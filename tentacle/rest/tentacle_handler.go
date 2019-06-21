package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/AHerczeg/ostopus/shared"
	"github.com/AHerczeg/ostopus/tentacle/local"
	tentacleQuery "github.com/AHerczeg/ostopus/tentacle/query"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func StartServing(address string) {
	MustStartRouter(address, setupRouter)
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
		WriteResponse(w, http.StatusBadRequest, []byte("error parsing head address"))
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

	client := getDefaultClient()
	resp, err := client.Post(headAddress, "application/json", bytes.NewBuffer(marshaledSelf))
	if err != nil {
		logrus.Error(err)
		WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while registering tentacle"))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		local.SetHeadAddress(headAddress)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	WriteResponse(w, resp.StatusCode, respBody)
}

func receiveCommand(w http.ResponseWriter, r *http.Request) {
	var query shared.Query
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		logrus.Error(err)
		WriteResponse(w, http.StatusBadRequest, []byte("failed to unmarshal query"))
		return
	}

	if query.Frequency != "" {
		// TODO handle frequency
	} else {
		if !query.Validate() {
			WriteResponse(w, http.StatusBadRequest, []byte("malformed query or frequency"))
			return
		}

		queryHandler := tentacleQuery.GetQueryHandler()

		result, err := queryHandler.RunCustomQuery(query.Query)
		fmt.Println(result)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, []byte(fmt.Sprintf("failed to run query. err: %v", err)))
			return
		}

		marshaledResults, err := json.Marshal(result)
		fmt.Println(string(marshaledResults))
		if err != nil {
			logrus.Error(err)
			WriteResponse(w, http.StatusInternalServerError, []byte("unexpected error while parsing result"))
		}

		WriteResponse(w, http.StatusOK, marshaledResults)
	}
}

func ping(w http.ResponseWriter, _ *http.Request) {
	logrus.Info("Receiving ping")
	WriteResponse(w, http.StatusOK, []byte{})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	WriteResponse(w, http.StatusNotFound, []byte{})
}

func WriteResponse(w http.ResponseWriter, code int, response []byte) {
	w.WriteHeader(code)
	if len(response) > 0 {
		w.Write(response)
	}
}

func getDefaultClient() http.Client {
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

	//router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
