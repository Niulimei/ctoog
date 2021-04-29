package restapi

import (
	"os/exec"
)

func InitRepoHandler() {

}

func Exec(commandLine string) ([]byte, error) {
	cmd := exec.Command("/bin/bash", "-c", commandLine)
	return cmd.Output()
}
