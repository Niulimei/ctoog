// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// ListPvobComponentStreamOKCode is the HTTP code returned for type ListPvobComponentStreamOK
const ListPvobComponentStreamOKCode int = 200

/*ListPvobComponentStreamOK 流列表

swagger:response listPvobComponentStreamOK
*/
type ListPvobComponentStreamOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewListPvobComponentStreamOK creates ListPvobComponentStreamOK with default headers values
func NewListPvobComponentStreamOK() *ListPvobComponentStreamOK {

	return &ListPvobComponentStreamOK{}
}

// WithPayload adds the payload to the list pvob component stream o k response
func (o *ListPvobComponentStreamOK) WithPayload(payload []string) *ListPvobComponentStreamOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list pvob component stream o k response
func (o *ListPvobComponentStreamOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPvobComponentStreamOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

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

// ListPvobComponentStreamInternalServerErrorCode is the HTTP code returned for type ListPvobComponentStreamInternalServerError
const ListPvobComponentStreamInternalServerErrorCode int = 500

/*ListPvobComponentStreamInternalServerError 内部错误

swagger:response listPvobComponentStreamInternalServerError
*/
type ListPvobComponentStreamInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewListPvobComponentStreamInternalServerError creates ListPvobComponentStreamInternalServerError with default headers values
func NewListPvobComponentStreamInternalServerError() *ListPvobComponentStreamInternalServerError {

	return &ListPvobComponentStreamInternalServerError{}
}

// WithPayload adds the payload to the list pvob component stream internal server error response
func (o *ListPvobComponentStreamInternalServerError) WithPayload(payload *models.ErrorModel) *ListPvobComponentStreamInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list pvob component stream internal server error response
func (o *ListPvobComponentStreamInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPvobComponentStreamInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
