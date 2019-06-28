// This file is safe to edit. Once it exists it will not be overwritten

package rest

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AHerczeg/ostopus/head/api/model"
	"github.com/AHerczeg/ostopus/head/tentacles"
	"github.com/AHerczeg/ostopus/shared"
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

type restResponse struct {
	code  int
	error error
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

	if api.TentaclePingTentaclesHandler == nil {
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
	}
	if api.TentacleQueryTentaclesHandler == nil {
		api.TentacleQueryTentaclesHandler = tentacle.QueryTentaclesHandlerFunc(func(params tentacle.QueryTentaclesParams) middleware.Responder {
			return middleware.NotImplemented("operation tentacle.QueryTentacles has not yet been implemented")
		})
	}
	if api.TentacleRegisterTentacleHandler == nil {
		api.TentacleRegisterTentacleHandler = tentacle.RegisterTentacleHandlerFunc(func(params tentacle.RegisterTentacleParams) middleware.Responder {
			return middleware.NotImplemented("operation tentacle.RegisterTentacle has not yet been implemented")
		})
	}
	if api.TentacleRemoveTentacleHandler == nil {
		api.TentacleRemoveTentacleHandler = tentacle.RemoveTentacleHandlerFunc(func(params tentacle.RemoveTentacleParams) middleware.Responder {
			return middleware.NotImplemented("operation tentacle.RemoveTentacle has not yet been implemented")
		})
	}

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
