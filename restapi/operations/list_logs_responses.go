// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// ListLogsCreatedCode is the HTTP code returned for type ListLogsCreated
const ListLogsCreatedCode int = 201

/*ListLogsCreated 成功

swagger:response listLogsCreated
*/
type ListLogsCreated struct {

	/*
	  In: Body
	*/
	Payload *models.LogPageInfoModel `json:"body,omitempty"`
}

// NewListLogsCreated creates ListLogsCreated with default headers values
func NewListLogsCreated() *ListLogsCreated {

	return &ListLogsCreated{}
}

// WithPayload adds the payload to the list logs created response
func (o *ListLogsCreated) WithPayload(payload *models.LogPageInfoModel) *ListLogsCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list logs created response
func (o *ListLogsCreated) SetPayload(payload *models.LogPageInfoModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListLogsCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListLogsInternalServerErrorCode is the HTTP code returned for type ListLogsInternalServerError
const ListLogsInternalServerErrorCode int = 500

/*ListLogsInternalServerError 内部错误

swagger:response listLogsInternalServerError
*/
type ListLogsInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewListLogsInternalServerError creates ListLogsInternalServerError with default headers values
func NewListLogsInternalServerError() *ListLogsInternalServerError {

	return &ListLogsInternalServerError{}
}

// WithPayload adds the payload to the list logs internal server error response
func (o *ListLogsInternalServerError) WithPayload(payload *models.ErrorModel) *ListLogsInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list logs internal server error response
func (o *ListLogsInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListLogsInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
