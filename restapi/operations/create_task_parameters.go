// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	"ctgb/models"
)

// NewCreateTaskParams creates a new CreateTaskParams object
//
// There are no default values defined in the spec.
func NewCreateTaskParams() CreateTaskParams {

	return CreateTaskParams{}
}

// CreateTaskParams contains all the bound params for the create task operation
// typically these are obtained from a http.Request
//
// swagger:parameters CreateTask
type CreateTaskParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*任务信息
	  Required: true
	  In: body
	*/
	UserInfo *models.TaskModel
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewCreateTaskParams() beforehand.
func (o *CreateTaskParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.TaskModel
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("userInfo", "body", ""))
			} else {
				res = append(res, errors.NewParseError("userInfo", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.UserInfo = &body
			}
		}
	} else {
		res = append(res, errors.Required("userInfo", "body", ""))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
