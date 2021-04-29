// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ListPvobComponentHandlerFunc turns a function with the right signature into a list pvob component handler
type ListPvobComponentHandlerFunc func(ListPvobComponentParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListPvobComponentHandlerFunc) Handle(params ListPvobComponentParams) middleware.Responder {
	return fn(params)
}

// ListPvobComponentHandler interface for that can handle valid list pvob component params
type ListPvobComponentHandler interface {
	Handle(ListPvobComponentParams) middleware.Responder
}

// NewListPvobComponent creates a new http.Handler for the list pvob component operation
func NewListPvobComponent(ctx *middleware.Context, handler ListPvobComponentHandler) *ListPvobComponent {
	return &ListPvobComponent{Context: ctx, Handler: handler}
}

/* ListPvobComponent swagger:route GET /pvobs/{id}/components listPvobComponent

组件列表

*/
type ListPvobComponent struct {
	Context *middleware.Context
	Handler ListPvobComponentHandler
}

func (o *ListPvobComponent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListPvobComponentParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}