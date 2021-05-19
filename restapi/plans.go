package restapi

import (
	"ctgb/database"
	"ctgb/models"
	"ctgb/restapi/operations"
	"ctgb/utils"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
)

var planColumns = []string{"id", "status", "origin_type", "pvob", "component", "dir", "origin_url", "translate_type", "target_url", "subsystem", "config_lib", "business_group", "team", "supporter", "supporter_tel", "creator", "tip", "project_type", "purpose", "plan_start_time", "plan_switch_time", "actual_start_time", "actual_switch_time", "effect", "extra1", "extra2", "extra3"}

func buildParams(params operations.ListPlanParams) map[string]string {
	//TODO
	return map[string]string{}
}
func buildPlanWhereSQL(queryParams map[string]string) (string, []interface{}, error) {
	l := len(queryParams)
	if l > 0 {
		sqlKeys := make([]string, 0, l)
		sqlValues := make([]interface{}, 0, l)

		placeholderIndex := int32(1)
		for k, v := range queryParams {
			switch k {
			case "status":
				sqlKeys, sqlValues, placeholderIndex = utils.GeneWhereLike(k, v, placeholderIndex, sqlKeys, sqlValues)
			}
		}
		return utils.GeneWhereSQL(sqlKeys, sqlValues)
	}
	return "", nil, nil
}

func ListPlanHandler(params operations.ListPlanParams) middleware.Responder {
	username, verified := utils.Verify(params.Authorization)
	if !verified {
		return middleware.Error(http.StatusUnauthorized, models.ErrorModel{Message: "鉴权失败"})
	}
	whereSQL, _, sqlErr := buildPlanWhereSQL(buildParams(params))
	if nil != sqlErr {
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: ""})
	}

	user := getUserInfo(username)
	if user.RoleID != int64(AdminRole) {
		whereSQL += fmt.Sprintf(" and creator=%s", username)
	}
	prepSQL := utils.PreparingQurySQL(planColumns, "plan", int(params.Offset), int(params.Limit), "id DESC", whereSQL)

	var plans []*database.PlanModel
	var count int64
	err := database.DB.Select(&plans, prepSQL, params.Limit, params.Offset)
	if err != nil {
		log.Error(err)
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: "Sql Error"})
	}
	queryCount := "select count(id) from plan where 1=1 " + whereSQL
	err = database.DB.Get(&count, queryCount)
	if err != nil {
		log.Error(err)
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: "Sql Error"})
	}
	plansPage := &models.PlanPageInfoModel{}
	var planModel []*models.PlanModel
	err = copier.Copy(&planModel, &plans)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: err.Error()})
	}
	plansPage.PlanInfo = planModel
	plansPage.Count = count
	return operations.NewListPlanOK().WithPayload(plansPage)
}

func CreatePlanHandler(params operations.CreatePlanParams) middleware.Responder {
	username, verified := utils.Verify(params.Authorization)
	if !verified {
		return middleware.Error(http.StatusUnauthorized, models.ErrorModel{Message: "鉴权失败"})
	}
	var plan database.PlanModel

	err := copier.Copy(&plan, params.PlanInfo)
	if err != nil {
		return middleware.Error(http.StatusInternalServerError, models.ErrorModel{Message: ""})
	}
	plan.Creator = username
	var ph []string
	for i := 1; i <= len(planColumns[1:]); i++ {
		ph = append(ph, "?")
	}
	sqlStr := fmt.Sprintf("INSERT INTO plan (%s) VALUES (%s)",
		strings.Join(planColumns[1:], ","), strings.Join(ph, ","))
	_, err = database.DB.Exec(sqlStr,
		plan.Status,
		plan.OriginType,
		plan.Pvob,
		plan.Component,
		plan.Dir,
		plan.OriginURL,
		plan.TranslateType,
		plan.TargetURL,
		plan.Subsystem,
		plan.ConfigLib,
		plan.Group,
		plan.Team,
		plan.Supporter,
		plan.SupporterTel,
		plan.Creator,
		plan.Tip,
		plan.ProjectType,
		plan.Purpose,
		plan.PlanStartTime,
		plan.PlanSwitchTime,
		plan.ActualStartTime,
		plan.ActualSwitchTime,
		plan.Effect,
		plan.Extra1,
		plan.Extra2,
		plan.Extra3,
	)
	if err != nil {
		return operations.NewDeletePlanInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	return operations.NewDeletePlanCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}

func GetPlanHandler(params operations.GetPlanParams) middleware.Responder {
	var plan = &database.PlanModel{}
	err := database.DB.Get(plan, "select * from plan where id=?", params.ID)
	if err != nil && err != sql.ErrNoRows {
		return operations.NewGetPlanInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	var planModel models.PlanModel
	err = copier.Copy(&planModel, plan)
	if err != nil {
		return operations.NewGetPlanInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return operations.NewGetPlanOK().WithPayload(&planModel)
}

func DeletePlanHandler(params operations.DeletePlanParams) middleware.Responder {
	_, err := database.DB.Exec("delete from plan where id=?", params.ID)
	if err != nil {
		return operations.NewDeletePlanInternalServerError().WithPayload(&models.ErrorModel{
			Code:    http.StatusInternalServerError,
			Message: "Sql Error",
		})
	}
	return operations.NewDeletePlanCreated().WithPayload(&models.OK{
		Message: "ok",
	})
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
			"UPDATE plan SET origin_type = :origin_type, pvob = :pvob, component = :component, dir = :dir,"+
				"origin_url = :origin_url, translate_type = :translate_type, target_url = :target_url, "+
				"subsystem = :subsystem, config_lib = :config_lib, business_group = :group, team = :team, supporter = :supporter,"+
				"supporter_tel = :supporter_tel, tip = :tip, project_type = :project_type, "+
				"purpose = :purpose, plan_start_time = :plan_start_time, plan_switch_time = :plan_switch_time, "+
				"effect = :effect WHERE id = :id", plan,
		)
	}
	return operations.NewUpdatePlanCreated().WithPayload(&models.OK{
		Message: "ok",
	})
}
