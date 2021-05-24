// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// ListWorkersOKCode is the HTTP code returned for type ListWorkersOK
const ListWorkersOKCode int = 200

/*ListWorkersOK 工作节点列表

swagger:response listWorkersOK
*/
type ListWorkersOK struct {

	/*
	  In: Body
	*/
	Payload *models.WorkerPageInfoModel `json:"body,omitempty"`
}

// NewListWorkersOK creates ListWorkersOK with default headers values
func NewListWorkersOK() *ListWorkersOK {

	return &ListWorkersOK{}
}

// WithPayload adds the payload to the list workers o k response
func (o *ListWorkersOK) WithPayload(payload *models.WorkerPageInfoModel) *ListWorkersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list workers o k response
func (o *ListWorkersOK) SetPayload(payload *models.WorkerPageInfoModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListWorkersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListWorkersInternalServerErrorCode is the HTTP code returned for type ListWorkersInternalServerError
const ListWorkersInternalServerErrorCode int = 500

/*ListWorkersInternalServerError 内部错误

swagger:response listWorkersInternalServerError
*/
type ListWorkersInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewListWorkersInternalServerError creates ListWorkersInternalServerError with default headers values
func NewListWorkersInternalServerError() *ListWorkersInternalServerError {

	return &ListWorkersInternalServerError{}
}

// WithPayload adds the payload to the list workers internal server error response
func (o *ListWorkersInternalServerError) WithPayload(payload *models.ErrorModel) *ListWorkersInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list workers internal server error response
func (o *ListWorkersInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListWorkersInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}