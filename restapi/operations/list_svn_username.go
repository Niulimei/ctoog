// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// ListSvnUsernameHandlerFunc turns a function with the right signature into a list svn username handler
type ListSvnUsernameHandlerFunc func(ListSvnUsernameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListSvnUsernameHandlerFunc) Handle(params ListSvnUsernameParams) middleware.Responder {
	return fn(params)
}

// ListSvnUsernameHandler interface for that can handle valid list svn username params
type ListSvnUsernameHandler interface {
	Handle(ListSvnUsernameParams) middleware.Responder
}

// NewListSvnUsername creates a new http.Handler for the list svn username operation
func NewListSvnUsername(ctx *middleware.Context, handler ListSvnUsernameHandler) *ListSvnUsername {
	return &ListSvnUsername{Context: ctx, Handler: handler}
}

/* ListSvnUsername swagger:route GET /svn_username_pairs/{id} listSvnUsername

svn用户名列表列表

*/
type ListSvnUsername struct {
	Context *middleware.Context
	Handler ListSvnUsernameHandler
}

func (o *ListSvnUsername) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListSvnUsernameParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
