// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ListWorkersHandlerFunc turns a function with the right signature into a list workers handler
type ListWorkersHandlerFunc func(ListWorkersParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListWorkersHandlerFunc) Handle(params ListWorkersParams) middleware.Responder {
	return fn(params)
}

// ListWorkersHandler interface for that can handle valid list workers params
type ListWorkersHandler interface {
	Handle(ListWorkersParams) middleware.Responder
}

// NewListWorkers creates a new http.Handler for the list workers operation
func NewListWorkers(ctx *middleware.Context, handler ListWorkersHandler) *ListWorkers {
	return &ListWorkers{Context: ctx, Handler: handler}
}

/* ListWorkers swagger:route GET /workers listWorkers

工作节点列表

*/
type ListWorkers struct {
	Context *middleware.Context
	Handler ListWorkersHandler
}

func (o *ListWorkers) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewListWorkersParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}