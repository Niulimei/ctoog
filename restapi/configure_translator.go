// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"ctgb/utils"
	"fmt"
	"net/http"
	"strings"

	"ctgb/restapi/operations"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
)

//go:generate swagger generate server --target ../../ctgb --name Translator --spec ../api/backend.yml --principal interface{}

func configureFlags(api *operations.TranslatorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.TranslatorAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.CreateUserHandler = operations.CreateUserHandlerFunc(CreateUserHandler)
	api.ListUserHandler = operations.ListUserHandlerFunc(ListUsersHandler)
	api.LoginHandler = operations.LoginHandlerFunc(LoginHandler)
	api.CreateTaskHandler = operations.CreateTaskHandlerFunc(CreateTaskHandler)
	api.GetTaskHandler = operations.GetTaskHandlerFunc(GetTaskHandler)
	api.ListTaskHandler = operations.ListTaskHandlerFunc(ListTaskHandler)
	api.UpdateTaskHandler = operations.UpdateTaskHandlerFunc(UpdateTaskHandler)
	api.RestartTaskHandler = operations.RestartTaskHandlerFunc(RestartTaskHandler)
	api.PingWorkerHandler = operations.PingWorkerHandlerFunc(PingWorkerHandler)
	api.ListPvobHandler = operations.ListPvobHandlerFunc(ListPvobHandler)
	api.ListPvobComponentHandler = operations.ListPvobComponentHandlerFunc(ListPvobComponentHandler)
	api.ListPvobComponentStreamHandler = operations.ListPvobComponentStreamHandlerFunc(ListPvobComponentStreamHandler)
	api.GetUserHandler = operations.GetUserHandlerFunc(GetUserHandler)
	api.ListLogsHandler = operations.ListLogsHandlerFunc(ListLogsHandler)
	api.GetTaskCommandOutHandler = operations.GetTaskCommandOutHandlerFunc(GetTaskCommandOutHandler)
	api.UpdateTaskCommandOutHandler = operations.UpdateTaskCommandOutHandlerFunc(UpdateTaskCommandOutHandler)
	api.DeleteTaskHandler = operations.DeleteTaskHandlerFunc(DeleteTaskHandler)
	api.DeleteTaskCacheHandler = operations.DeleteTaskCacheHandlerFunc(DeleteTaskCacheHandler)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	go utils.LogHandle()

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		username, valid := utils.Verify(token)
		if !valid && !strings.HasSuffix(r.RequestURI, "/login") {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		r.Header.Set("username", username)
		handler.ServeHTTP(w, r)
	})
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	defer func() {
		if ret := recover(); ret != nil {
			fmt.Printf("Recover From Panic. %v\n", ret)
		}
	}()
	return handler
}
