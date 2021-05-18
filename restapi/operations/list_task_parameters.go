// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewListTaskParams creates a new ListTaskParams object
// with the default values initialized.
func NewListTaskParams() ListTaskParams {

	var (
		// initialize parameters with default values

		limitDefault  = int64(0)
		offsetDefault = int64(0)
	)

	return ListTaskParams{
		Limit: limitDefault,

		Offset: offsetDefault,
	}
}

// ListTaskParams contains all the bound params for the list task operation
// typically these are obtained from a http.Request
//
// swagger:parameters ListTask
type ListTaskParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: header
	*/
	Authorization string
	/*
	  In: query
	*/
	Component *string
	/*
	  Required: true
	  In: query
	  Default: 0
	*/
	Limit int64
	/*
	  Required: true
	  In: query
	  Default: 0
	*/
	Offset int64
	/*
	  In: query
	*/
	Pvob *string
	/*
	  In: query
	*/
	Stream *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewListTaskParams() beforehand.
func (o *ListTaskParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if err := o.bindAuthorization(r.Header[http.CanonicalHeaderKey("Authorization")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	qComponent, qhkComponent, _ := qs.GetOK("component")
	if err := o.bindComponent(qComponent, qhkComponent, route.Formats); err != nil {
		res = append(res, err)
	}

	qLimit, qhkLimit, _ := qs.GetOK("limit")
	if err := o.bindLimit(qLimit, qhkLimit, route.Formats); err != nil {
		res = append(res, err)
	}

	qOffset, qhkOffset, _ := qs.GetOK("offset")
	if err := o.bindOffset(qOffset, qhkOffset, route.Formats); err != nil {
		res = append(res, err)
	}

	qPvob, qhkPvob, _ := qs.GetOK("pvob")
	if err := o.bindPvob(qPvob, qhkPvob, route.Formats); err != nil {
		res = append(res, err)
	}

	qStream, qhkStream, _ := qs.GetOK("stream")
	if err := o.bindStream(qStream, qhkStream, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuthorization binds and validates parameter Authorization from header.
func (o *ListTaskParams) bindAuthorization(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindComponent binds and validates parameter Component from query.
func (o *ListTaskParams) bindComponent(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Component = &raw

	return nil
}

// bindLimit binds and validates parameter Limit from query.
func (o *ListTaskParams) bindLimit(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("limit", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("limit", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("limit", "query", "int64", raw)
	}
	o.Limit = value

	return nil
}

// bindOffset binds and validates parameter Offset from query.
func (o *ListTaskParams) bindOffset(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("offset", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("offset", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("offset", "query", "int64", raw)
	}
	o.Offset = value

	return nil
}

// bindPvob binds and validates parameter Pvob from query.
func (o *ListTaskParams) bindPvob(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Pvob = &raw

	return nil
}

// bindStream binds and validates parameter Stream from query.
func (o *ListTaskParams) bindStream(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Stream = &raw

	return nil
}
