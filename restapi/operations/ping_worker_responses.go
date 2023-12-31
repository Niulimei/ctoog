// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// PingWorkerCreatedCode is the HTTP code returned for type PingWorkerCreated
const PingWorkerCreatedCode int = 201

/*PingWorkerCreated 成功

swagger:response pingWorkerCreated
*/
type PingWorkerCreated struct {

	/*
	  In: Body
	*/
	Payload *models.OK `json:"body,omitempty"`
}

// NewPingWorkerCreated creates PingWorkerCreated with default headers values
func NewPingWorkerCreated() *PingWorkerCreated {

	return &PingWorkerCreated{}
}

// WithPayload adds the payload to the ping worker created response
func (o *PingWorkerCreated) WithPayload(payload *models.OK) *PingWorkerCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the ping worker created response
func (o *PingWorkerCreated) SetPayload(payload *models.OK) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PingWorkerCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PingWorkerInternalServerErrorCode is the HTTP code returned for type PingWorkerInternalServerError
const PingWorkerInternalServerErrorCode int = 500

/*PingWorkerInternalServerError 内部错误

swagger:response pingWorkerInternalServerError
*/
type PingWorkerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewPingWorkerInternalServerError creates PingWorkerInternalServerError with default headers values
func NewPingWorkerInternalServerError() *PingWorkerInternalServerError {

	return &PingWorkerInternalServerError{}
}

// WithPayload adds the payload to the ping worker internal server error response
func (o *PingWorkerInternalServerError) WithPayload(payload *models.ErrorModel) *PingWorkerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the ping worker internal server error response
func (o *PingWorkerInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PingWorkerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
