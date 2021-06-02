package restapi

import (
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

var configs = []string{"svnRoute"}

//var configs = []string{"svnRoute"}

func GetFrontendConfig(params operations.GetFrontConfigParams) middleware.Responder {
	return operations.NewGetFrontConfigOK().WithPayload([]string{})
}
