// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// GetUserOKCode is the HTTP code returned for type GetUserOK
const GetUserOKCode int = 200

/*GetUserOK 账户

swagger:response getUserOK
*/
type GetUserOK struct {

	/*
	  In: Body
	*/
	Payload *models.UserInfoModel `json:"body,omitempty"`
}

// NewGetUserOK creates GetUserOK with default headers values
func NewGetUserOK() *GetUserOK {

	return &GetUserOK{}
}

// WithPayload adds the payload to the get user o k response
func (o *GetUserOK) WithPayload(payload *models.UserInfoModel) *GetUserOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user o k response
func (o *GetUserOK) SetPayload(payload *models.UserInfoModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetUserInternalServerErrorCode is the HTTP code returned for type GetUserInternalServerError
const GetUserInternalServerErrorCode int = 500

/*GetUserInternalServerError 内部错误

swagger:response getUserInternalServerError
*/
type GetUserInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewGetUserInternalServerError creates GetUserInternalServerError with default headers values
func NewGetUserInternalServerError() *GetUserInternalServerError {

	return &GetUserInternalServerError{}
}

// WithPayload adds the payload to the get user internal server error response
func (o *GetUserInternalServerError) WithPayload(payload *models.ErrorModel) *GetUserInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get user internal server error response
func (o *GetUserInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUserInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}