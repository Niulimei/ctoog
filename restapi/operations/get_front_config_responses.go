// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// GetFrontConfigOKCode is the HTTP code returned for type GetFrontConfigOK
const GetFrontConfigOKCode int = 200

/*GetFrontConfigOK 前端配置

swagger:response getFrontConfigOK
*/
type GetFrontConfigOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewGetFrontConfigOK creates GetFrontConfigOK with default headers values
func NewGetFrontConfigOK() *GetFrontConfigOK {

	return &GetFrontConfigOK{}
}

// WithPayload adds the payload to the get front config o k response
func (o *GetFrontConfigOK) WithPayload(payload []string) *GetFrontConfigOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get front config o k response
func (o *GetFrontConfigOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetFrontConfigOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]string, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetFrontConfigInternalServerErrorCode is the HTTP code returned for type GetFrontConfigInternalServerError
const GetFrontConfigInternalServerErrorCode int = 500

/*GetFrontConfigInternalServerError 内部错误

swagger:response getFrontConfigInternalServerError
*/
type GetFrontConfigInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewGetFrontConfigInternalServerError creates GetFrontConfigInternalServerError with default headers values
func NewGetFrontConfigInternalServerError() *GetFrontConfigInternalServerError {

	return &GetFrontConfigInternalServerError{}
}

// WithPayload adds the payload to the get front config internal server error response
func (o *GetFrontConfigInternalServerError) WithPayload(payload *models.ErrorModel) *GetFrontConfigInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get front config internal server error response
func (o *GetFrontConfigInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetFrontConfigInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}