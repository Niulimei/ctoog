// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// UpdateTaskCreatedCode is the HTTP code returned for type UpdateTaskCreated
const UpdateTaskCreatedCode int = 201

/*UpdateTaskCreated 成功

swagger:response updateTaskCreated
*/
type UpdateTaskCreated struct {

	/*
	  In: Body
	*/
	Payload *models.OK `json:"body,omitempty"`
}

// NewUpdateTaskCreated creates UpdateTaskCreated with default headers values
func NewUpdateTaskCreated() *UpdateTaskCreated {

	return &UpdateTaskCreated{}
}

// WithPayload adds the payload to the update task created response
func (o *UpdateTaskCreated) WithPayload(payload *models.OK) *UpdateTaskCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update task created response
func (o *UpdateTaskCreated) SetPayload(payload *models.OK) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTaskCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTaskInternalServerErrorCode is the HTTP code returned for type UpdateTaskInternalServerError
const UpdateTaskInternalServerErrorCode int = 500

/*UpdateTaskInternalServerError 内部错误

swagger:response updateTaskInternalServerError
*/
type UpdateTaskInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewUpdateTaskInternalServerError creates UpdateTaskInternalServerError with default headers values
func NewUpdateTaskInternalServerError() *UpdateTaskInternalServerError {

	return &UpdateTaskInternalServerError{}
}

// WithPayload adds the payload to the update task internal server error response
func (o *UpdateTaskInternalServerError) WithPayload(payload *models.ErrorModel) *UpdateTaskInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update task internal server error response
func (o *UpdateTaskInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTaskInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
