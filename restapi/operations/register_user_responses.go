// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// RegisterUserCreatedCode is the HTTP code returned for type RegisterUserCreated
const RegisterUserCreatedCode int = 201

/*RegisterUserCreated 成功

swagger:response registerUserCreated
*/
type RegisterUserCreated struct {

	/*
	  In: Body
	*/
	Payload *models.OK `json:"body,omitempty"`
}

// NewRegisterUserCreated creates RegisterUserCreated with default headers values
func NewRegisterUserCreated() *RegisterUserCreated {

	return &RegisterUserCreated{}
}

// WithPayload adds the payload to the register user created response
func (o *RegisterUserCreated) WithPayload(payload *models.OK) *RegisterUserCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the register user created response
func (o *RegisterUserCreated) SetPayload(payload *models.OK) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RegisterUserCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// RegisterUserInternalServerErrorCode is the HTTP code returned for type RegisterUserInternalServerError
const RegisterUserInternalServerErrorCode int = 500

/*RegisterUserInternalServerError 内部错误

swagger:response registerUserInternalServerError
*/
type RegisterUserInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewRegisterUserInternalServerError creates RegisterUserInternalServerError with default headers values
func NewRegisterUserInternalServerError() *RegisterUserInternalServerError {

	return &RegisterUserInternalServerError{}
}

// WithPayload adds the payload to the register user internal server error response
func (o *RegisterUserInternalServerError) WithPayload(payload *models.ErrorModel) *RegisterUserInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the register user internal server error response
func (o *RegisterUserInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *RegisterUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
