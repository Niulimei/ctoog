// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// WorkerDetail worker detail
//
// swagger:model WorkerDetail
type WorkerDetail struct {

	// id
	ID int64 `json:"id,omitempty"`

	// register time
	RegisterTime string `json:"registerTime,omitempty" db:"register_time"`

	// status
	Status string `json:"status,omitempty"`

	// task count
	TaskCount int64 `json:"taskCount,omitempty" db:"task_count"`

	// worker Url
	WorkerURL string `json:"workerUrl,omitempty" db:"worker_url"`
}

// Validate validates this worker detail
func (m *WorkerDetail) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this worker detail based on context it is used
func (m *WorkerDetail) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *WorkerDetail) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkerDetail) UnmarshalBinary(b []byte) error {
	var res WorkerDetail
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}