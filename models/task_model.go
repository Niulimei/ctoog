// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type JsonNullInt64 struct {
	sql.NullInt64
}

func (v JsonNullInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int64)
	} else {
		return json.Marshal(0)
	}
}

func (v *JsonNullInt64) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Int64 = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullString struct {
	sql.NullString
}

func (v JsonNullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal("")
	}
}

func (v *JsonNullString) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}

type JsonNullBool struct {
	sql.NullBool
}

func (v JsonNullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	} else {
		return json.Marshal(false)
	}
}

func (v *JsonNullBool) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *bool
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Bool = *x
	} else {
		v.Valid = false
	}
	return nil
}

// TaskModel task model
//
// swagger:model TaskModel
type TaskModel struct {

	// branches info
	BranchesInfo string `json:"branches_info,omitempty" db:"branches_info"`

	// cc password
	CcPassword JsonNullString `json:"ccPassword,omitempty" db:"cc_password"`

	// cc user
	CcUser JsonNullString `json:"ccUser,omitempty" db:"cc_user"`

	// component
	Component JsonNullString `json:"component,omitempty"`

	// dir
	Dir JsonNullString `json:"dir,omitempty"`

	// git email
	GitEmail JsonNullString `json:"gitEmail,omitempty" db:"git_email"`

	// git password
	GitPassword JsonNullString `json:"gitPassword,omitempty" db:"git_password"`

	// git URL
	GitURL JsonNullString `json:"gitURL,omitempty" db:"git_url"`

	// git user
	GitUser JsonNullString `json:"gitUser,omitempty" db:"git_user"`

	// gitee group
	GiteeGroup JsonNullString `json:"giteeGroup,omitempty" db:"gitee_group"`

	// gitee project
	GiteeProject JsonNullString `json:"giteeProject,omitempty" db:"gitee_project"`

	// gitee token
	GiteeToken JsonNullString `json:"giteeToken,omitempty" db:"gitee_token"`

	// gitignore
	Gitignore JsonNullString `json:"gitignore,omitempty"`

	// gitlab group
	GitlabGroup JsonNullString `json:"gitlabGroup,omitempty" db:"gitlab_group"`

	// gitlab project
	GitlabProject JsonNullString `json:"gitlabProject,omitempty" db:"gitlab_project"`

	// gitlab token
	GitlabToken JsonNullString `json:"gitlabToken,omitempty" db:"gitlab_token"`

	// gitlab source
	SourceURL JsonNullString `json:"sourceURL,omitempty" db:"source_url"`

	// gitlab target
	TargetURL JsonNullString `json:"targetURL,omitempty" db:"target_url"`

	// include empty
	IncludeEmpty JsonNullBool `json:"includeEmpty,omitempty" db:"include_empty"`

	// keep
	Keep JsonNullString `json:"keep,omitempty"`

	// match info
	MatchInfo []*TaskMatchInfo `json:"matchInfo"`

	// model type
	ModelType JsonNullString `json:"modelType,omitempty" db:"model_type"`

	// name pair
	NamePair []*NamePairInfo `json:"namePair"`

	// pvob
	Pvob JsonNullString `json:"pvob,omitempty"`

	// status
	Status JsonNullString `json:"status,omitempty"`

	// svn Url
	SvnURL JsonNullString `json:"svn_url,omitempty" db:"svn_url"`

	// worker Url
	WorkerURL JsonNullString `json:"workerUrl,omitempty"`

	// worker id
	WorkerId JsonNullInt64 `json:"workerId,omitempty" db:"worker_id"`
}

// Validate validates this task model
func (m *TaskModel) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMatchInfo(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNamePair(formats); err != nil {
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

func (m *TaskModel) validateNamePair(formats strfmt.Registry) error {
	if swag.IsZero(m.NamePair) { // not required
		return nil
	}

	for i := 0; i < len(m.NamePair); i++ {
		if swag.IsZero(m.NamePair[i]) { // not required
			continue
		}

		if m.NamePair[i] != nil {
			if err := m.NamePair[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("namePair" + "." + strconv.Itoa(i))
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

	if err := m.contextValidateNamePair(ctx, formats); err != nil {
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

func (m *TaskModel) contextValidateNamePair(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.NamePair); i++ {

		if m.NamePair[i] != nil {
			if err := m.NamePair[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("namePair" + "." + strconv.Itoa(i))
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
