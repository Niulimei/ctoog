// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewTranslatorAPI creates a new Translator instance
func NewTranslatorAPI(spec *loads.Document) *TranslatorAPI {
	return &TranslatorAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),

		JSONProducer: runtime.JSONProducer(),

		CreateTaskHandler: CreateTaskHandlerFunc(func(params CreateTaskParams) middleware.Responder {
			return middleware.NotImplemented("operation CreateTask has not yet been implemented")
		}),
		CreateUserHandler: CreateUserHandlerFunc(func(params CreateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation CreateUser has not yet been implemented")
		}),
		GetTaskHandler: GetTaskHandlerFunc(func(params GetTaskParams) middleware.Responder {
			return middleware.NotImplemented("operation GetTask has not yet been implemented")
		}),
		GetUserHandler: GetUserHandlerFunc(func(params GetUserParams) middleware.Responder {
			return middleware.NotImplemented("operation GetUser has not yet been implemented")
		}),
		ListPvobHandler: ListPvobHandlerFunc(func(params ListPvobParams) middleware.Responder {
			return middleware.NotImplemented("operation ListPvob has not yet been implemented")
		}),
		ListPvobComponentHandler: ListPvobComponentHandlerFunc(func(params ListPvobComponentParams) middleware.Responder {
			return middleware.NotImplemented("operation ListPvobComponent has not yet been implemented")
		}),
		ListPvobComponentStreamHandler: ListPvobComponentStreamHandlerFunc(func(params ListPvobComponentStreamParams) middleware.Responder {
			return middleware.NotImplemented("operation ListPvobComponentStream has not yet been implemented")
		}),
		ListTaskHandler: ListTaskHandlerFunc(func(params ListTaskParams) middleware.Responder {
			return middleware.NotImplemented("operation ListTask has not yet been implemented")
		}),
		ListUserHandler: ListUserHandlerFunc(func(params ListUserParams) middleware.Responder {
			return middleware.NotImplemented("operation ListUser has not yet been implemented")
		}),
		LoginHandler: LoginHandlerFunc(func(params LoginParams) middleware.Responder {
			return middleware.NotImplemented("operation Login has not yet been implemented")
		}),
		PingWorkerHandler: PingWorkerHandlerFunc(func(params PingWorkerParams) middleware.Responder {
			return middleware.NotImplemented("operation PingWorker has not yet been implemented")
		}),
		RestartTaskHandler: RestartTaskHandlerFunc(func(params RestartTaskParams) middleware.Responder {
			return middleware.NotImplemented("operation RestartTask has not yet been implemented")
		}),
		UpdateTaskHandler: UpdateTaskHandlerFunc(func(params UpdateTaskParams) middleware.Responder {
			return middleware.NotImplemented("operation UpdateTask has not yet been implemented")
		}),
	}
}

/*TranslatorAPI the translator API */
type TranslatorAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer

	// CreateTaskHandler sets the operation handler for the create task operation
	CreateTaskHandler CreateTaskHandler
	// CreateUserHandler sets the operation handler for the create user operation
	CreateUserHandler CreateUserHandler
	// GetTaskHandler sets the operation handler for the get task operation
	GetTaskHandler GetTaskHandler
	// GetUserHandler sets the operation handler for the get user operation
	GetUserHandler GetUserHandler
	// ListPvobHandler sets the operation handler for the list pvob operation
	ListPvobHandler ListPvobHandler
	// ListPvobComponentHandler sets the operation handler for the list pvob component operation
	ListPvobComponentHandler ListPvobComponentHandler
	// ListPvobComponentStreamHandler sets the operation handler for the list pvob component stream operation
	ListPvobComponentStreamHandler ListPvobComponentStreamHandler
	// ListTaskHandler sets the operation handler for the list task operation
	ListTaskHandler ListTaskHandler
	// ListUserHandler sets the operation handler for the list user operation
	ListUserHandler ListUserHandler
	// LoginHandler sets the operation handler for the login operation
	LoginHandler LoginHandler
	// PingWorkerHandler sets the operation handler for the ping worker operation
	PingWorkerHandler PingWorkerHandler
	// RestartTaskHandler sets the operation handler for the restart task operation
	RestartTaskHandler RestartTaskHandler
	// UpdateTaskHandler sets the operation handler for the update task operation
	UpdateTaskHandler UpdateTaskHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *TranslatorAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *TranslatorAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *TranslatorAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *TranslatorAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *TranslatorAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *TranslatorAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *TranslatorAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *TranslatorAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *TranslatorAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the TranslatorAPI
func (o *TranslatorAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.CreateTaskHandler == nil {
		unregistered = append(unregistered, "CreateTaskHandler")
	}
	if o.CreateUserHandler == nil {
		unregistered = append(unregistered, "CreateUserHandler")
	}
	if o.GetTaskHandler == nil {
		unregistered = append(unregistered, "GetTaskHandler")
	}
	if o.GetUserHandler == nil {
		unregistered = append(unregistered, "GetUserHandler")
	}
	if o.ListPvobHandler == nil {
		unregistered = append(unregistered, "ListPvobHandler")
	}
	if o.ListPvobComponentHandler == nil {
		unregistered = append(unregistered, "ListPvobComponentHandler")
	}
	if o.ListPvobComponentStreamHandler == nil {
		unregistered = append(unregistered, "ListPvobComponentStreamHandler")
	}
	if o.ListTaskHandler == nil {
		unregistered = append(unregistered, "ListTaskHandler")
	}
	if o.ListUserHandler == nil {
		unregistered = append(unregistered, "ListUserHandler")
	}
	if o.LoginHandler == nil {
		unregistered = append(unregistered, "LoginHandler")
	}
	if o.PingWorkerHandler == nil {
		unregistered = append(unregistered, "PingWorkerHandler")
	}
	if o.RestartTaskHandler == nil {
		unregistered = append(unregistered, "RestartTaskHandler")
	}
	if o.UpdateTaskHandler == nil {
		unregistered = append(unregistered, "UpdateTaskHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *TranslatorAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *TranslatorAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *TranslatorAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *TranslatorAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *TranslatorAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *TranslatorAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the translator API
func (o *TranslatorAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *TranslatorAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/tasks"] = NewCreateTask(o.context, o.CreateTaskHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/users"] = NewCreateUser(o.context, o.CreateUserHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/tasks/{id}"] = NewGetTask(o.context, o.GetTaskHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users/self"] = NewGetUser(o.context, o.GetUserHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pvobs"] = NewListPvob(o.context, o.ListPvobHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pvobs/{id}/components"] = NewListPvobComponent(o.context, o.ListPvobComponentHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/pvobs/{pvob_id}/components/{component_id}"] = NewListPvobComponentStream(o.context, o.ListPvobComponentStreamHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/tasks"] = NewListTask(o.context, o.ListTaskHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/users"] = NewListUser(o.context, o.ListUserHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/login"] = NewLogin(o.context, o.LoginHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/workers"] = NewPingWorker(o.context, o.PingWorkerHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/tasks/restart"] = NewRestartTask(o.context, o.RestartTaskHandler)
	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/tasks/{id}"] = NewUpdateTask(o.context, o.UpdateTaskHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *TranslatorAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *TranslatorAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *TranslatorAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *TranslatorAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *TranslatorAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
