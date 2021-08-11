// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CCHistoryInfoModel c c history info model
//
// swagger:model CCHistoryInfoModel
type CCHistoryInfoModel struct {

	// count
	Count int64 `json:"count,omitempty"`

	// create time
	CreateTime string `json:"createTime,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// history type
	HistoryType string `json:"historyType,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// key
	Key int64 `json:"key,omitempty"`

	// limit
	Limit int64 `json:"limit,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// offset
	Offset int64 `json:"offset,omitempty"`

	// owner
	Owner string `json:"owner,omitempty"`
}

// Validate validates this c c history info model
func (m *CCHistoryInfoModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this c c history info model based on context it is used
func (m *CCHistoryInfoModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CCHistoryInfoModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CCHistoryInfoModel) UnmarshalBinary(b []byte) error {
	var res CCHistoryInfoModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
