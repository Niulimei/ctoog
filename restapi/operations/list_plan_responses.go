// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// ListPlanOKCode is the HTTP code returned for type ListPlanOK
const ListPlanOKCode int = 200

/*ListPlanOK 迁移计划列表

swagger:response listPlanOK
*/
type ListPlanOK struct {

	/*
	  In: Body
	*/
	Payload *models.PlanPageInfoModel `json:"body,omitempty"`
}

// NewListPlanOK creates ListPlanOK with default headers values
func NewListPlanOK() *ListPlanOK {

	return &ListPlanOK{}
}

// WithPayload adds the payload to the list plan o k response
func (o *ListPlanOK) WithPayload(payload *models.PlanPageInfoModel) *ListPlanOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list plan o k response
func (o *ListPlanOK) SetPayload(payload *models.PlanPageInfoModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPlanOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListPlanInternalServerErrorCode is the HTTP code returned for type ListPlanInternalServerError
const ListPlanInternalServerErrorCode int = 500

/*ListPlanInternalServerError 内部错误

swagger:response listPlanInternalServerError
*/
type ListPlanInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewListPlanInternalServerError creates ListPlanInternalServerError with default headers values
func NewListPlanInternalServerError() *ListPlanInternalServerError {

	return &ListPlanInternalServerError{}
}

// WithPayload adds the payload to the list plan internal server error response
func (o *ListPlanInternalServerError) WithPayload(payload *models.ErrorModel) *ListPlanInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list plan internal server error response
func (o *ListPlanInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListPlanInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
