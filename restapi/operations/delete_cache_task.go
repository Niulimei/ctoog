// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteCacheTaskHandlerFunc turns a function with the right signature into a delete cache task handler
type DeleteCacheTaskHandlerFunc func(DeleteCacheTaskParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteCacheTaskHandlerFunc) Handle(params DeleteCacheTaskParams) middleware.Responder {
	return fn(params)
}

// DeleteCacheTaskHandler interface for that can handle valid delete cache task params
type DeleteCacheTaskHandler interface {
	Handle(DeleteCacheTaskParams) middleware.Responder
}

// NewDeleteCacheTask creates a new http.Handler for the delete cache task operation
func NewDeleteCacheTask(ctx *middleware.Context, handler DeleteCacheTaskHandler) *DeleteCacheTask {
	return &DeleteCacheTask{Context: ctx, Handler: handler}
}

/* DeleteCacheTask swagger:route DELETE /tasks/cache/{id} deleteCacheTask

任务缓存删除

*/
type DeleteCacheTask struct {
	Context *middleware.Context
	Handler DeleteCacheTaskHandler
}

func (o *DeleteCacheTask) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteCacheTaskParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}