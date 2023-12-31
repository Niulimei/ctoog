// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// GetTaskCommandOutHandlerFunc turns a function with the right signature into a get task command out handler
type GetTaskCommandOutHandlerFunc func(GetTaskCommandOutParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetTaskCommandOutHandlerFunc) Handle(params GetTaskCommandOutParams) middleware.Responder {
	return fn(params)
}

// GetTaskCommandOutHandler interface for that can handle valid get task command out params
type GetTaskCommandOutHandler interface {
	Handle(GetTaskCommandOutParams) middleware.Responder
}

// NewGetTaskCommandOut creates a new http.Handler for the get task command out operation
func NewGetTaskCommandOut(ctx *middleware.Context, handler GetTaskCommandOutHandler) *GetTaskCommandOut {
	return &GetTaskCommandOut{Context: ctx, Handler: handler}
}

/* GetTaskCommandOut swagger:route GET /tasks/cmdout/{log_id} getTaskCommandOut

任务执行详情

*/
type GetTaskCommandOut struct {
	Context *middleware.Context
	Handler GetTaskCommandOutHandler
}

func (o *GetTaskCommandOut) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetTaskCommandOutParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
