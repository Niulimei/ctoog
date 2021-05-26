// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"

	"ctgb/models"
)

// NewUpdateSvnNamePairParams creates a new UpdateSvnNamePairParams object
//
// There are no default values defined in the spec.
func NewUpdateSvnNamePairParams() UpdateSvnNamePairParams {

	return UpdateSvnNamePairParams{}
}

// UpdateSvnNamePairParams contains all the bound params for the update svn name pair operation
// typically these are obtained from a http.Request
//
// swagger:parameters UpdateSvnNamePair
type UpdateSvnNamePairParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	ID string
	/*对应信息
	  Required: true
	  In: body
	*/
	UsernamePairInfo []*models.NamePairInfo
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUpdateSvnNamePairParams() beforehand.
func (o *UpdateSvnNamePairParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body []*models.NamePairInfo
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("usernamePairInfo", "body", ""))
			} else {
				res = append(res, errors.NewParseError("usernamePairInfo", "body", "", err))
			}
		} else {

			// validate array of body objects
			for i := range body {
				if body[i] == nil {
					continue
				}
				if err := body[i].Validate(route.Formats); err != nil {
					res = append(res, err)
					break
				}
			}

			if len(res) == 0 {
				o.UsernamePairInfo = body
			}
		}
	} else {
		res = append(res, errors.Required("usernamePairInfo", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindID binds and validates parameter ID from path.
func (o *UpdateSvnNamePairParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.ID = raw

	return nil
}