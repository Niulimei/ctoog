// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TaskRestart task restart
//
// swagger:model TaskRestart
type TaskRestart struct {

	// id
	ID string `json:"id,omitempty"`
}

// Validate validates this task restart
func (m *TaskRestart) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this task restart based on context it is used
func (m *TaskRestart) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TaskRestart) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskRestart) UnmarshalBinary(b []byte) error {
	var res TaskRestart
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}