// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteTaskCacheHandlerFunc turns a function with the right signature into a delete task cache handler
type DeleteTaskCacheHandlerFunc func(DeleteTaskCacheParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteTaskCacheHandlerFunc) Handle(params DeleteTaskCacheParams) middleware.Responder {
	return fn(params)
}

// DeleteTaskCacheHandler interface for that can handle valid delete task cache params
type DeleteTaskCacheHandler interface {
	Handle(DeleteTaskCacheParams) middleware.Responder
}

// NewDeleteTaskCache creates a new http.Handler for the delete task cache operation
func NewDeleteTaskCache(ctx *middleware.Context, handler DeleteTaskCacheHandler) *DeleteTaskCache {
	return &DeleteTaskCache{Context: ctx, Handler: handler}
}

/* DeleteTaskCache swagger:route DELETE /tasks/cache/{id} deleteTaskCache

任务缓存删除

*/
type DeleteTaskCache struct {
	Context *middleware.Context
	Handler DeleteTaskCacheHandler
}

func (o *DeleteTaskCache) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteTaskCacheParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
