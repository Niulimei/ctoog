package restapi

import (
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
)

func GetFrontendConfig(params operations.GetFrontConfigParams) middleware.Responder {
	return operations.NewGetFrontConfigOK().WithPayload([]string{"svnRoute"})
}
