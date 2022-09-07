// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"ctgb/models"
)

// GetCCHistoryOKCode is the HTTP code returned for type GetCCHistoryOK
const GetCCHistoryOKCode int = 200

/*GetCCHistoryOK cc历史信息

swagger:response getCCHistoryOK
*/
type GetCCHistoryOK struct {

	/*
	  In: Body
	*/
	Payload []*models.CCHistoryInfoModel `json:"body,omitempty"`
}

// NewGetCCHistoryOK creates GetCCHistoryOK with default headers values
func NewGetCCHistoryOK() *GetCCHistoryOK {

	return &GetCCHistoryOK{}
}

// WithPayload adds the payload to the get c c history o k response
func (o *GetCCHistoryOK) WithPayload(payload []*models.CCHistoryInfoModel) *GetCCHistoryOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get c c history o k response
func (o *GetCCHistoryOK) SetPayload(payload []*models.CCHistoryInfoModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCCHistoryOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*models.CCHistoryInfoModel, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

// GetCCHistoryInternalServerErrorCode is the HTTP code returned for type GetCCHistoryInternalServerError
const GetCCHistoryInternalServerErrorCode int = 500

/*GetCCHistoryInternalServerError 内部错误

swagger:response getCCHistoryInternalServerError
*/
type GetCCHistoryInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorModel `json:"body,omitempty"`
}

// NewGetCCHistoryInternalServerError creates GetCCHistoryInternalServerError with default headers values
func NewGetCCHistoryInternalServerError() *GetCCHistoryInternalServerError {

	return &GetCCHistoryInternalServerError{}
}

// WithPayload adds the payload to the get c c history internal server error response
func (o *GetCCHistoryInternalServerError) WithPayload(payload *models.ErrorModel) *GetCCHistoryInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get c c history internal server error response
func (o *GetCCHistoryInternalServerError) SetPayload(payload *models.ErrorModel) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetCCHistoryInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}