// This file is safe to edit. Once it exists it will not be overwritten

package rest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/AHerczeg/ostopus/head/api/model"
	"github.com/AHerczeg/ostopus/head/tentacles"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/AHerczeg/ostopus/head/api/rest/operation"
	"github.com/AHerczeg/ostopus/head/api/rest/operation/tentacle"
)

//go:generate swagger generate server --target ../../api --name Ostopus --spec ../swagger.yml --api-package operation --model-package model --server-package rest --exclude-main

type pingResponse struct {
	tentacle string
	response bool
}

func configureFlags(api *operation.OstopusAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operation.OstopusAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TentacleRegisterTentacleHandler = tentacle.RegisterTentacleHandlerFunc(func(params tentacle.RegisterTentacleParams) middleware.Responder {
		id := registerTentacle(params.Body)
		if id == 0 {
			return tentacle.NewPingTentaclesDefault(500)
		}
		return tentacle.NewRegisterTentacleCreated().WithPayload(id)
	})

	api.TentaclePingTentaclesHandler = tentacle.PingTentaclesHandlerFunc(func(params tentacle.PingTentaclesParams) middleware.Responder {
		responses, err := pingTentacles()
		if err != nil {
			switch err.Code {
			case 500:
				return tentacle.NewPingTentaclesInternalServerError().WithPayload(err)
			default:
				return tentacle.NewPingTentaclesDefault(500)
			}
		}
		return tentacle.NewPingTentaclesOK().WithPayload(responses)
	})

	api.TentacleQueryTentaclesHandler = tentacle.QueryTentaclesHandlerFunc(func(params tentacle.QueryTentaclesParams) middleware.Responder {
		result, err := relayQuery(params.Body)
		if err != nil {
			switch err.Code {
			case 500:
				return tentacle.NewQueryTentaclesInternalServerError().WithPayload(err)
			default:
				return tentacle.NewQueryTentaclesDefault(500)
			}
		}
		return tentacle.NewQueryTentaclesOK().WithPayload(result)
	})

	api.TentacleRemoveTentacleHandler = tentacle.RemoveTentacleHandlerFunc(func(params tentacle.RemoveTentacleParams) middleware.Responder {
		if err := removeTentacle(params.ID); err != nil {
			switch err.Code {
			case 404:
				return tentacle.NewRemoveTentacleNotFound()
			default:
				return tentacle.NewRemoveTentacleDefault(500)
			}
		}
		return tentacle.NewRemoveTentacleDefault(204)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}

func registerTentacle(tentacle *model.Tentacle) int64 {
	tentacle.ID = time.Now().UnixNano()
	tentacles.Tentacles().SaveTentacle(*tentacle)
	logrus.WithFields(logrus.Fields{"Name": tentacle.Name, "Address": tentacle.Address}).Info("New tentacle registered")
	return tentacle.ID
}

func pingTentacles() (string, *model.Error) {
	results := make(map[string]bool)
	allTentacles := tentacles.Tentacles().GetAllTentacles()

	// If there are no tentacles we can quit early
	if len(allTentacles) == 0 {
		return "{}", nil
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
		return "", &model.Error{Code: 500, Message: "unexpected error while preparing response"}
	}

	return string(marshaledResults), nil
}

func pingTentacle(tentacle model.Tentacle, response chan<- pingResponse) {
	logrus.WithFields(logrus.Fields{"name": tentacle.Name, "address": tentacle.Address}).Info("Pinging tentacle")
	res, err := http.Get(*tentacle.Address + "/ping")
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

func relayQuery(query *model.Query) (*model.Result, *model.Error) {

	//TODO validate query

	marshaledCommand, err := json.Marshal(query.Command)
	if err != nil {
		logrus.Error(err)
		return nil, &model.Error{Code: 500, Message: "unexpected error while parsing query"}
	}

	var results sync.Map
	var wg sync.WaitGroup

	for _, target := range query.Target {
		if storedTentacle, ok := tentacles.Tentacles().GetTentacle(target); ok {
			wg.Add(1)
			go syncQuery(marshaledCommand, &storedTentacle, &results, &wg)
		} else {
			results.Store(target, "unknown target")
		}
	}

	wg.Wait()

	return mapToResult(results), nil

}

func syncQuery(query []byte, tentacle *model.Tentacle, results *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	results.Store(tentacle.ID, sendQuery(query, *tentacle.Address))
}

func sendQuery(query []byte, address string) string {
	req, err := http.NewRequest("POST", address+"/query", bytes.NewBuffer(query))

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return "error while processing response"
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("failed query: code %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "error while reading response"
	}
	return string(bodyBytes)

}

func mapToResult(syncMap sync.Map) *model.Result {
	var result model.Result

	syncMap.Range(func(k, v interface{}) bool {
		result.Payload = append(result.Payload, &model.ResultPayloadItems0{ID: k.(int64), Result: v.(string)})
		return true
	})
	return &result
}

func removeTentacle(id int64) *model.Error {
	if ok := tentacles.Tentacles().RemoveTentacle(id); !ok {
		return &model.Error{Code: 404, Message: "tentacle not found"}
	}
	return nil
}
