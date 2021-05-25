// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ListPlanHandlerFunc turns a function with the right signature into a list plan handler
type ListPlanHandlerFunc func(ListPlanParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListPlanHandlerFunc) Handle(params ListPlanParams) middleware.Responder {
	return fn(params)
}

// ListPlanHandler interface for that can handle valid list plan params
type ListPlanHandler interface {
	Handle(ListPlanParams) middleware.Responder
}

// NewListPlan creates a new http.Handler for the list plan operation
func NewListPlan(ctx *middleware.Context, handler ListPlanHandler) *ListPlan {
	return &ListPlan{Context: ctx, Handler: handler}
}

/* ListPlan swagger:route GET /plans listPlan

迁移计划列表

*/
type ListPlan struct {
	Context *middleware.Context
	Handler ListPlanHandler
}

func (o *ListPlan) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewListPlanParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
