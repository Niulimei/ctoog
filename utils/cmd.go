package utils

import (
	"os/exec"
)

func Exec(commandLine string) ([]byte, error) {
	cmd := exec.Command("/bin/bash", "-c", commandLine)
	return cmd.Output()
}
