package restapi

import (
	"ctgb/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"os/exec"
	"strings"
)

func GetAllPvob() []string {
	pvobs := make([]string, 10)
	cmd := exec.Command("cleartool", "lsvob", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	result := string(out)
	infos := strings.Split(result, "\n\n")
	for _, info := range infos {
		lines := strings.Split(info, "\n")
		if len(lines) == 0 {
			continue
		}
		if strings.HasPrefix(lines[0], "Tag: ") && lines[len(lines)-1] == "Vob registry attributes: ucmvob" {
			pvob := lines[0][5:]
			pvobs = append(pvobs, pvob)
		}
	}
	return pvobs
}

func GetAllComponent(pvob string) []string {
	components := make([]string, 10)
	args := `lscomp -fmt "%[root_dir]p\n" -invob ` + pvob
	cmd := exec.Command("cleartool", strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	result := string(out)
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			components = append(components, line)
		}
	}
	return components
}

func checkStreamComponent(pvob, component, stream string) bool {
	args := `lsstream -fmt %[components]p ` + stream + " " + pvob
	cmd := exec.Command("cleartool", strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	result := string(out)
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if line == component {
			return true
		}
	}
	return false
}

func GetAllStream(pvob, component string) []string {
	streams := make([]string, 10)
	args := `lsstream -s -invob ` + pvob
	cmd := exec.Command("cleartool", strings.Split(args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	result := string(out)
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if checkStreamComponent(pvob, component, line) {
			streams = append(streams, line)
		}
	}
	return streams
}

func ListPvobHandler(params operations.ListPvobParams) middleware.Responder {
	return operations.NewListPvobOK().WithPayload(GetAllPvob())
}

func ListPvobComponentHandler(params operations.ListPvobComponentParams) middleware.Responder {
	pvob := params.ID
	return operations.NewListPvobComponentOK().WithPayload(GetAllComponent(pvob))
}

func ListPvobComponentStreamHandler(params operations.ListPvobComponentStreamParams) middleware.Responder {
	pvob := params.PvobID
	component := params.ComponentID
	return operations.NewListPvobComponentStreamOK().WithPayload(GetAllStream(pvob, component))
}
