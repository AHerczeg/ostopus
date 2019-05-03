package shared

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func WriteResponse(w http.ResponseWriter, code int, response []byte) {
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
	if err := StartRouter(address, routerSetup); err != nil {
		panic(err)
	}
}

func StartRouter(address string, routerSetup func(*mux.Router)) error {
	logrus.Info("Starting up router")
	router := mux.NewRouter()
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
