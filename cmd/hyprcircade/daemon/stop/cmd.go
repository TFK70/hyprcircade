package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/tfk70/hyprcircade/internal/logging"

	"github.com/urfave/cli/v3"
)

func CreateDaemonStopCmd(bin string) *cli.Command {
	run := func(context context.Context, cmd *cli.Command) error {
		logger, err := logging.GetLogger()
		if err != nil {
			return err
		}

		pidFile := fmt.Sprintf("/tmp/%s.pid", bin)

		data, err := os.ReadFile(pidFile)
		if err != nil {
			return fmt.Errorf("Unable to read pid file. Is daemon running?\n%v", err)
		}

		pidStr := strings.TrimSpace(string(data))
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return err
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}

		err = proc.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}

		logger.Info("Daemon stopped")

		return nil
	}

	return &cli.Command{
		Name:   "stop",
		Usage:  "stop hyprcircade daemon (if it is running in background mode)",
		Action: run,
	}
}
