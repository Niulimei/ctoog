// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PlanModel plan model
//
// swagger:model PlanModel
type PlanModel struct {

	// actual start time
	ActualStartTime string `json:"actual_start_time,omitempty"`

	// actual switch time
	ActualSwitchTime string `json:"actual_switch_time,omitempty"`

	// component
	Component string `json:"component,omitempty"`

	// 配置库
	ConfigLib string `json:"configLib,omitempty"`

	// dir
	Dir string `json:"dir,omitempty"`

	// 影响范围
	Effect string `json:"effect,omitempty"`

	// 事业群
	Group string `json:"group,omitempty"`

	// 计划编号
	ID int64 `json:"id,omitempty"`

	// 源仓库类型（cc、gerrit、私服）
	OriginType string `json:"originType,omitempty"`

	// gerrit或者私服的git地址
	OriginURL string `json:"originUrl,omitempty"`

	// plan start time
	PlanStartTime string `json:"plan_start_time,omitempty"`

	// plan switch time
	PlanSwitchTime string `json:"plan_switch_time,omitempty"`

	// 工程类型
	ProjectType string `json:"projectType,omitempty"`

	// 业务用途
	Purpose string `json:"purpose,omitempty"`

	// pvob
	Pvob string `json:"pvob,omitempty"`

	// 状态
	Status string `json:"status,omitempty"`

	// 物理子系统
	Subsystem string `json:"subsystem,omitempty"`

	// 对接人姓名
	Supporter string `json:"supporter,omitempty"`

	// 对接人电话
	SupporterTel string `json:"supporterTel,omitempty"`

	// 目标git地址
	TargetURL string `json:"targetUrl,omitempty"`

	// task id
	TaskID int64 `json:"task_id,omitempty"`

	// 项目组
	Team string `json:"team,omitempty"`

	// 备注
	Tip string `json:"tip,omitempty"`

	// 自己迁移还是工作组迁移
	TranslateType string `json:"translateType,omitempty"`
}

// Validate validates this plan model
func (m *PlanModel) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this plan model based on context it is used
func (m *PlanModel) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *PlanModel) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PlanModel) UnmarshalBinary(b []byte) error {
	var res PlanModel
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
