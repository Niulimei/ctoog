// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateTaskCommandOutHandlerFunc turns a function with the right signature into a update task command out handler
type UpdateTaskCommandOutHandlerFunc func(UpdateTaskCommandOutParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateTaskCommandOutHandlerFunc) Handle(params UpdateTaskCommandOutParams) middleware.Responder {
	return fn(params)
}

// UpdateTaskCommandOutHandler interface for that can handle valid update task command out params
type UpdateTaskCommandOutHandler interface {
	Handle(UpdateTaskCommandOutParams) middleware.Responder
}

// NewUpdateTaskCommandOut creates a new http.Handler for the update task command out operation
func NewUpdateTaskCommandOut(ctx *middleware.Context, handler UpdateTaskCommandOutHandler) *UpdateTaskCommandOut {
	return &UpdateTaskCommandOut{Context: ctx, Handler: handler}
}

/* UpdateTaskCommandOut swagger:route POST /tasks/cmdout/{log_id} updateTaskCommandOut

更新任务执行详情

*/
type UpdateTaskCommandOut struct {
	Context *middleware.Context
	Handler UpdateTaskCommandOutHandler
}

func (o *UpdateTaskCommandOut) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateTaskCommandOutParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}