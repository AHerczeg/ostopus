// This file is safe to edit. Once it exists it will not be overwritten

package rest

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/AHerczeg/ostopus/head/api/rest/operation"
	"github.com/AHerczeg/ostopus/head/api/rest/operation/tentacle"
)

//go:generate swagger generate server --target ../../api --name Ostopus --spec ../swagger.yml --api-package operation --model-package model --server-package rest --exclude-main

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
			return middleware.NotImplemented("operation tentacle.PingTentacles has not yet been implemented")
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
