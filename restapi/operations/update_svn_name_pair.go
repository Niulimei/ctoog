// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// UpdateSvnNamePairHandlerFunc turns a function with the right signature into a update svn name pair handler
type UpdateSvnNamePairHandlerFunc func(UpdateSvnNamePairParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateSvnNamePairHandlerFunc) Handle(params UpdateSvnNamePairParams) middleware.Responder {
	return fn(params)
}

// UpdateSvnNamePairHandler interface for that can handle valid update svn name pair params
type UpdateSvnNamePairHandler interface {
	Handle(UpdateSvnNamePairParams) middleware.Responder
}

// NewUpdateSvnNamePair creates a new http.Handler for the update svn name pair operation
func NewUpdateSvnNamePair(ctx *middleware.Context, handler UpdateSvnNamePairHandler) *UpdateSvnNamePair {
	return &UpdateSvnNamePair{Context: ctx, Handler: handler}
}

/* UpdateSvnNamePair swagger:route PUT /svn_username_pairs/{id} updateSvnNamePair

更新用户名对应信息

*/
type UpdateSvnNamePair struct {
	Context *middleware.Context
	Handler UpdateSvnNamePairHandler
}

func (o *UpdateSvnNamePair) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUpdateSvnNamePairParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
