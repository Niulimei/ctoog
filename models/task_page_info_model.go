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

// TaskPageInfoModel task page info model
//
// swagger:model TaskPageInfoModel
type TaskPageInfoModel struct {

	// count
	Count int64 `json:"count,omitempty"`

	// limit
	Limit int64 `json:"limit,omitempty"`

	// offset
	Offset int64 `json:"offset,omitempty"`

	// task info
	TaskInfo []*TaskInfoModel `json:"taskInfo"`
}

// Validate validates this task page info model
func (m *TaskPageInfoModel) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTaskInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskPageInfoModel) validateTaskInfo(formats strfmt.Registry) error {
	if swag.IsZero(m.TaskInfo) { // not required
		return nil
	}

	for i := 0; i < len(m.TaskInfo); i++ {
		if swag.IsZero(m.TaskInfo[i]) { // not required
			continue
		}

		if m.TaskInfo[i] != nil {
			if err := m.TaskInfo[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("taskInfo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this task page info model based on the context it is used
func (m *TaskPageInfoModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTaskInfo(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TaskPageInfoModel) contextValidateTaskInfo(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.TaskInfo); i++ {

		if m.TaskInfo[i] != nil {
			if err := m.TaskInfo[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("taskInfo" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *TaskPageInfoModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TaskPageInfoModel) UnmarshalBinary(b []byte) error {
	var res TaskPageInfoModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
