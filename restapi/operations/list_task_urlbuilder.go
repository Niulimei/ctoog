// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// ListTaskURL generates an URL for the list task operation
type ListTaskURL struct {
	Component *string
	Limit     int64
	ModelType *string
	Offset    int64
	Pvob      *string
	Status    *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ListTaskURL) WithBasePath(bp string) *ListTaskURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *ListTaskURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *ListTaskURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/tasks"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var componentQ string
	if o.Component != nil {
		componentQ = *o.Component
	}
	if componentQ != "" {
		qs.Set("component", componentQ)
	}

	limitQ := swag.FormatInt64(o.Limit)
	if limitQ != "" {
		qs.Set("limit", limitQ)
	}

	var modelTypeQ string
	if o.ModelType != nil {
		modelTypeQ = *o.ModelType
	}
	if modelTypeQ != "" {
		qs.Set("modelType", modelTypeQ)
	}

	offsetQ := swag.FormatInt64(o.Offset)
	if offsetQ != "" {
		qs.Set("offset", offsetQ)
	}

	var pvobQ string
	if o.Pvob != nil {
		pvobQ = *o.Pvob
	}
	if pvobQ != "" {
		qs.Set("pvob", pvobQ)
	}

	var statusQ string
	if o.Status != nil {
		statusQ = *o.Status
	}
	if statusQ != "" {
		qs.Set("status", statusQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *ListTaskURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *ListTaskURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *ListTaskURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on ListTaskURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on ListTaskURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *ListTaskURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
