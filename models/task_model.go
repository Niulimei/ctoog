// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// TaskModel task model
//
// swagger:model TaskModel
type TaskModel struct {

	// cc password
	CcPassword string `json:"ccPassword,omitempty"`

	// cc user
	CcUser string `json:"ccUser,omitempty"`

	// component
	Component string `json:"component,omitempty"`

	// git password
	GitPassword string `json:"gitPassword,omitempty"`

	// git URL
	GitURL string `json:"gitURL,omitempty"`

	// git user
	GitUser string `json:"gitUser,omitempty"`

	// match info
	MatchInfo []*TaskMatchInfo `json:"matchInfo"`

	// pvob
	Pvob string `json:"pvob,omitempty"`
}

// Validate validates this task model
func (m *TaskModel) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMatchInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskModel) validateMatchInfo(formats strfmt.Registry) error {
	if swag.IsZero(m.MatchInfo) { // not required
		return nil
	}

	for i := 0; i < len(m.MatchInfo); i++ {
		if swag.IsZero(m.MatchInfo[i]) { // not required
			continue
		}

		if m.MatchInfo[i] != nil {
			if err := m.MatchInfo[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("matchInfo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this task model based on the context it is used
func (m *TaskModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMatchInfo(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskModel) contextValidateMatchInfo(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.MatchInfo); i++ {

		if m.MatchInfo[i] != nil {
			if err := m.MatchInfo[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("matchInfo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *TaskModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskModel) UnmarshalBinary(b []byte) error {
	var res TaskModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
