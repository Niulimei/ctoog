// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// DeletePlanCreatedCode is the HTTP code returned for type DeletePlanCreated
const DeletePlanCreatedCode int = 201

/*DeletePlanCreated 成功

swagger:response deletePlanCreated
*/
type DeletePlanCreated struct {

	/*
	  In: Body
	*/
	Payload *models.OK `json:"body,omitempty"`
}

// NewDeletePlanCreated creates DeletePlanCreated with default headers values
func NewDeletePlanCreated() *DeletePlanCreated {

	return &DeletePlanCreated{}
}

// WithPayload adds the payload to the delete plan created response
func (o *DeletePlanCreated) WithPayload(payload *models.OK) *DeletePlanCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete plan created response
func (o *DeletePlanCreated) SetPayload(payload *models.OK) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeletePlanCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeletePlanInternalServerErrorCode is the HTTP code returned for type DeletePlanInternalServerError
const DeletePlanInternalServerErrorCode int = 500

/*DeletePlanInternalServerError 内部错误

swagger:response deletePlanInternalServerError
*/
type DeletePlanInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewDeletePlanInternalServerError creates DeletePlanInternalServerError with default headers values
func NewDeletePlanInternalServerError() *DeletePlanInternalServerError {

	return &DeletePlanInternalServerError{}
}

// WithPayload adds the payload to the delete plan internal server error response
func (o *DeletePlanInternalServerError) WithPayload(payload *models.ErrorModel) *DeletePlanInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete plan internal server error response
func (o *DeletePlanInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeletePlanInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
