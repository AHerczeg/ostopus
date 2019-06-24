package rest

import (
	"net/http"
	"time"

	"github.com/AHerczeg/ostopus/head/rest/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func StartServing(address string) {
	MustStartRouter(address, setupRouter)
}

func setupRouter(router *mux.Router) {
	router.HandleFunc("/register", handlers.RegisterTentacle).Methods("POST")
	router.HandleFunc("/remove", handlers.RemoveTentacle).Methods("DELETE")

	router.HandleFunc("/ping", handlers.PingAll).Methods("GET")

	router.HandleFunc("/query", handlers.RelayQuery).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(handlers.NotFound)
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
