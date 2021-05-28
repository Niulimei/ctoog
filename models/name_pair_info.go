// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NamePairInfo name pair info
//
// swagger:model NamePairInfo
type NamePairInfo struct {

	// git email
	GitEmail string `json:"gitEmail,omitempty" db:"git_email"`

	// git user name
	GitUserName string `json:"gitUserName,omitempty" db:"git_username"`

	// svn user name
	SvnUserName string `json:"svnUserName,omitempty" db:"svn_username"`
}

// Validate validates this name pair info
func (m *NamePairInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this name pair info based on context it is used
func (m *NamePairInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *NamePairInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NamePairInfo) UnmarshalBinary(b []byte) error {
	var res NamePairInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
