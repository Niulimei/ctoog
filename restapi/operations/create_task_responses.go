// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// CreateTaskCreatedCode is the HTTP code returned for type CreateTaskCreated
const CreateTaskCreatedCode int = 201

/*CreateTaskCreated 成功

swagger:response createTaskCreated
*/
type CreateTaskCreated struct {

	/*
	  In: Body
	*/
	Payload *models.OK `json:"body,omitempty"`
}

// NewCreateTaskCreated creates CreateTaskCreated with default headers values
func NewCreateTaskCreated() *CreateTaskCreated {

	return &CreateTaskCreated{}
}

// WithPayload adds the payload to the create task created response
func (o *CreateTaskCreated) WithPayload(payload *models.OK) *CreateTaskCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create task created response
func (o *CreateTaskCreated) SetPayload(payload *models.OK) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTaskCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// CreateTaskInternalServerErrorCode is the HTTP code returned for type CreateTaskInternalServerError
const CreateTaskInternalServerErrorCode int = 500

/*CreateTaskInternalServerError 内部错误

swagger:response createTaskInternalServerError
*/
type CreateTaskInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewCreateTaskInternalServerError creates CreateTaskInternalServerError with default headers values
func NewCreateTaskInternalServerError() *CreateTaskInternalServerError {

	return &CreateTaskInternalServerError{}
}

// WithPayload adds the payload to the create task internal server error response
func (o *CreateTaskInternalServerError) WithPayload(payload *models.ErrorModel) *CreateTaskInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the create task internal server error response
func (o *CreateTaskInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CreateTaskInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
