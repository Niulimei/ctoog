// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetTaskCommandOutParams creates a new GetTaskCommandOutParams object
//
// There are no default values defined in the spec.
func NewGetTaskCommandOutParams() GetTaskCommandOutParams {

	return GetTaskCommandOutParams{}
}

// GetTaskCommandOutParams contains all the bound params for the get task command out operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetTaskCommandOut
type GetTaskCommandOutParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: header
	*/
	Authorization string
	/*
	  Required: true
	  In: path
	*/
	LogID int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetTaskCommandOutParams() beforehand.
func (o *GetTaskCommandOutParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindAuthorization(r.Header[http.CanonicalHeaderKey("Authorization")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	rLogID, rhkLogID, _ := route.Params.GetOK("log_id")
	if err := o.bindLogID(rLogID, rhkLogID, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuthorization binds and validates parameter Authorization from header.
func (o *GetTaskCommandOutParams) bindAuthorization(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("Authorization", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("Authorization", "header", raw); err != nil {
		return err
	}
	o.Authorization = raw

	return nil
}

// bindLogID binds and validates parameter LogID from path.
func (o *GetTaskCommandOutParams) bindLogID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("log_id", "path", "int64", raw)
	}
	o.LogID = value

	return nil
}
