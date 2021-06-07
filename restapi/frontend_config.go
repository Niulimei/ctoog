package restapi

import (
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"os"
)

func GetFrontendConfig(params operations.GetFrontConfigParams) middleware.Responder {
	var frontendConfigs []string
	_, ok := os.LookupEnv("SVN_SUPPORT")
	if ok {
		frontendConfigs = append(frontendConfigs, "svnRoute")
	}
	_, ok = os.LookupEnv("JIANXIN_SUPPORT")
	if ok {
		frontendConfigs = append(frontendConfigs, "jianxin")
	}
	return operations.NewGetFrontConfigOK().WithPayload(frontendConfigs)
}
