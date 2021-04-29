// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// ListPvobComponentOKCode is the HTTP code returned for type ListPvobComponentOK
const ListPvobComponentOKCode int = 200

/*ListPvobComponentOK 组件列表

swagger:response listPvobComponentOK
*/
type ListPvobComponentOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewListPvobComponentOK creates ListPvobComponentOK with default headers values
func NewListPvobComponentOK() *ListPvobComponentOK {

	return &ListPvobComponentOK{}
}

// WithPayload adds the payload to the list pvob component o k response
func (o *ListPvobComponentOK) WithPayload(payload []string) *ListPvobComponentOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list pvob component o k response
func (o *ListPvobComponentOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPvobComponentOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ListPvobComponentInternalServerErrorCode is the HTTP code returned for type ListPvobComponentInternalServerError
const ListPvobComponentInternalServerErrorCode int = 500

/*ListPvobComponentInternalServerError 内部错误

swagger:response listPvobComponentInternalServerError
*/
type ListPvobComponentInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewListPvobComponentInternalServerError creates ListPvobComponentInternalServerError with default headers values
func NewListPvobComponentInternalServerError() *ListPvobComponentInternalServerError {

	return &ListPvobComponentInternalServerError{}
}

// WithPayload adds the payload to the list pvob component internal server error response
func (o *ListPvobComponentInternalServerError) WithPayload(payload *models.ErrorModel) *ListPvobComponentInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list pvob component internal server error response
func (o *ListPvobComponentInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPvobComponentInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}