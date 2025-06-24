package commands

import (
	"fmt"
	"os/exec"

	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/sirupsen/logrus"
)

const (
	EXECUTABLE = "sh"
)

func ExecuteCommand(cmd string) error {
	logger, err := logging.GetLogger()
	if err != nil {
		return err
	}

	commandLogger := logger.WithFields(logrus.Fields{
		"cmd": cmd,
	})

	commandLogger.Info("Executing command")

	_, err = exec.LookPath(EXECUTABLE)
	if err != nil {
		return err
	}

	c := exec.Command(EXECUTABLE, "-c", cmd)

	output, err := c.CombinedOutput()
	if err != nil {
		formattedErr := fmt.Errorf("Command exited with error: %v\n%s", err, string(output))
		return formattedErr
	}

	commandLogger.Debugf("Cmd output:\n%s", string(output))
	commandLogger.Info("Successfully executed")

	return nil
}
