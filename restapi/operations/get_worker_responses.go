// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// GetWorkerOKCode is the HTTP code returned for type GetWorkerOK
const GetWorkerOKCode int = 200

/*GetWorkerOK 工作节点信息

swagger:response getWorkerOK
*/
type GetWorkerOK struct {

	/*
	  In: Body
	*/
	Payload *models.WorkerDetail `json:"body,omitempty"`
}

// NewGetWorkerOK creates GetWorkerOK with default headers values
func NewGetWorkerOK() *GetWorkerOK {

	return &GetWorkerOK{}
}

// WithPayload adds the payload to the get worker o k response
func (o *GetWorkerOK) WithPayload(payload *models.WorkerDetail) *GetWorkerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get worker o k response
func (o *GetWorkerOK) SetPayload(payload *models.WorkerDetail) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetWorkerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetWorkerInternalServerErrorCode is the HTTP code returned for type GetWorkerInternalServerError
const GetWorkerInternalServerErrorCode int = 500

/*GetWorkerInternalServerError 内部错误

swagger:response getWorkerInternalServerError
*/
type GetWorkerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewGetWorkerInternalServerError creates GetWorkerInternalServerError with default headers values
func NewGetWorkerInternalServerError() *GetWorkerInternalServerError {

	return &GetWorkerInternalServerError{}
}

// WithPayload adds the payload to the get worker internal server error response
func (o *GetWorkerInternalServerError) WithPayload(payload *models.ErrorModel) *GetWorkerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get worker internal server error response
func (o *GetWorkerInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetWorkerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
