// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// LogInfoModel log info model
//
// swagger:model LogInfoModel
type LogInfoModel struct {

	// action
	Action string `json:"action,omitempty"`

	// errcode
	Errcode int64 `json:"errcode,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// level
	Level string `json:"level,omitempty"`

	// message
	Message string `json:"message,omitempty"`

	// position
	Position string `json:"position,omitempty"`

	// time
	Time int64 `json:"time,omitempty"`

	// user
	User string `json:"user,omitempty"`
}

// Validate validates this log info model
func (m *LogInfoModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this log info model based on context it is used
func (m *LogInfoModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *LogInfoModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *LogInfoModel) UnmarshalBinary(b []byte) error {
	var res LogInfoModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
