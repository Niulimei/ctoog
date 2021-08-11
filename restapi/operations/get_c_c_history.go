// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetCCHistoryHandlerFunc turns a function with the right signature into a get c c history handler
type GetCCHistoryHandlerFunc func(GetCCHistoryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetCCHistoryHandlerFunc) Handle(params GetCCHistoryParams) middleware.Responder {
	return fn(params)
}

// GetCCHistoryHandler interface for that can handle valid get c c history params
type GetCCHistoryHandler interface {
	Handle(GetCCHistoryParams) middleware.Responder
}

// NewGetCCHistory creates a new http.Handler for the get c c history operation
func NewGetCCHistory(ctx *middleware.Context, handler GetCCHistoryHandler) *GetCCHistory {
	return &GetCCHistory{Context: ctx, Handler: handler}
}

/* GetCCHistory swagger:route GET /cc_history getCCHistory

cc历史信息

*/
type GetCCHistory struct {
	Context *middleware.Context
	Handler GetCCHistoryHandler
}

func (o *GetCCHistory) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewGetCCHistoryParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
