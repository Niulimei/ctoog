package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"github.com/go-openapi/runtime/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/jinzhu/copier"
)

func ListPlanHandler(params operations.ListPlanParams) middleware.Responder {
	return nil
}

func CreatePlanHandler(params operations.CreatePlanParams) middleware.Responder {
	return nil
}

func GetPlanHandler(params operations.GetPlanParams) middleware.Responder {
	return nil
}

func DeletePlanHandler(params operations.DeletePlanParams) middleware.Responder {
	return nil
}

func UpdatePlanHandler(params operations.UpdatePlanParams) middleware.Responder {
	planId := params.ID
	planParams := params.PlanInfo
	var plan database.PlanModel
	err := database.DB.Get(&plan, "SELECT * FROM plan WHERE id = $1", planId)
	if err != nil || plan.ID == 0 {
		log.Error("Get plan err:", err)
		return operations.NewUpdatePlanInternalServerError().WithPayload(&models.ErrorModel{
			Code:    404,
			Message: "没有发现计划",
		})
	}
	if planParams.Status != "" {
		if plan.TargetURL == "" {
			return operations.NewUpdatePlanInternalServerError().WithPayload(&models.ErrorModel{
				Code:    400,
				Message: "请填写git url后再执行操作",
			})
		}
		tx, _ := database.DB.Begin()
		if planParams.Status == "已迁移" {
			tx.Exec("UPDATE status = $1, actual_start_time = $2 WHERE id = $3",
				planParams.Status, planParams.ActualStartTime, planId)
		} else if planParams.Status == "已切换" {
			tx.Exec("UPDATE status = $1, actual_switch_time = $2 WHERE id = $3",
				planParams.Status, planParams.ActualSwitchTime, planId)
		} else if planParams.Status == "迁移中" {
			userToken := params.Authorization
			username, verified := utils.Verify(userToken)
			if !verified {
				return operations.NewUpdatePlanInternalServerError().WithPayload(&models.ErrorModel{
					Code:    401,
					Message: "鉴权失败",
				})
			}
			tx.Exec("UPDATE status = $1 WHERE id = $2",
				planParams.Status, planId)
			tx.Exec("INSERT INTO task (pvob, component, git_url,"+
				"status, last_completed_date_time, creator, dir, worker_id)"+
				" VALUES ($1, $2, $3, 'init', '', $4, $5, 0)",
				plan.Pvob, plan.Component, plan.TargetURL, username, plan.Dir)
		}
		tx.Commit()
	} else {
		err = copier.Copy(plan, planParams)
		if err != nil {
			log.Error(err)
			return operations.NewUpdatePlanInternalServerError().WithPayload(&models.ErrorModel{
				Code:    500,
				Message: err.Error(),
			})
		}
		database.DB.NamedExec(
			"UPDATE plan SET origin_type = :origin_type, pvob = :pvob, component = :component, dir = :dir," +
				"origin_url = :origin_url, translate_type = :translate_type, target_url = :target_url, " +
				"subsystem = :subsystem, config_lib = :config_lib, group = :group, team = :team, supporter = :supporter," +
				"supporter_tel = :supporter_tel, creator = :creator, tip = :tip, project_type = :project_type, " +
				"purpose = :purpose, plan_start_time = :plan_start_time, plan_switch_time = :plan_switch_time, " +
				"effect = :effect WHERE id = :id", plan,
		)
	}
	return operations.NewUpdatePlanCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
