package commands

import (
	"os/exec"
)

const (
	EXECUTABLE = "sh"
)

func ExecuteCommand(cmd string) error {
	_, err := exec.LookPath(EXECUTABLE)
	if err != nil {
		return err
	}

	c := exec.Command(EXECUTABLE, "-c", cmd)

	_, err = c.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
