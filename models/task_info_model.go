// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TaskInfoModel task info model
//
// swagger:model TaskInfoModel
type TaskInfoModel struct {

	// component
	Component string `json:"component,omitempty"`

	// dir
	Dir string `json:"dir,omitempty"`

	// git email
	GitEmail string `json:"gitEmail,omitempty" db:"git_email"`

	// git repo
	GitRepo string `json:"gitRepo,omitempty" db:"git_url"`

	// id
	ID int64 `json:"id,omitempty"`

	// include empty
	IncludeEmpty bool `json:"includeEmpty,omitempty" db:"include_empty"`

	// keep
	Keep string `json:"keep,omitempty"`

	// last complete date time
	LastCompleteDateTime string `json:"lastCompleteDateTime,omitempty" db:"last_completed_date_time"`

	// pvob
	Pvob string `json:"pvob,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// svn Url
	SvnURL string `json:"svnUrl,omitempty"`
}

// Validate validates this task info model
func (m *TaskInfoModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this task info model based on context it is used
func (m *TaskInfoModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *TaskInfoModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskInfoModel) UnmarshalBinary(b []byte) error {
	var res TaskInfoModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
