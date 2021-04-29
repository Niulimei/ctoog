// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// ListPvobOKCode is the HTTP code returned for type ListPvobOK
const ListPvobOKCode int = 200

/*ListPvobOK PVOB列表

swagger:response listPvobOK
*/
type ListPvobOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewListPvobOK creates ListPvobOK with default headers values
func NewListPvobOK() *ListPvobOK {

	return &ListPvobOK{}
}

// WithPayload adds the payload to the list pvob o k response
func (o *ListPvobOK) WithPayload(payload []string) *ListPvobOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list pvob o k response
func (o *ListPvobOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPvobOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ListPvobInternalServerErrorCode is the HTTP code returned for type ListPvobInternalServerError
const ListPvobInternalServerErrorCode int = 500

/*ListPvobInternalServerError 内部错误

swagger:response listPvobInternalServerError
*/
type ListPvobInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewListPvobInternalServerError creates ListPvobInternalServerError with default headers values
func NewListPvobInternalServerError() *ListPvobInternalServerError {

	return &ListPvobInternalServerError{}
}

// WithPayload adds the payload to the list pvob internal server error response
func (o *ListPvobInternalServerError) WithPayload(payload *models.ErrorModel) *ListPvobInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list pvob internal server error response
func (o *ListPvobInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPvobInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}