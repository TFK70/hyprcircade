package cmd

import (
	daemonLogs "github.com/tfk70/hyprcircade/cmd/hyprcircade/daemon/logs"
	daemonStop "github.com/tfk70/hyprcircade/cmd/hyprcircade/daemon/stop"
	"github.com/urfave/cli/v3"
)

func CreateDaemonCmd(bin string) *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			daemonStop.CreateDaemonStopCmd(bin),
			daemonLogs.CreateDaemonLogsCmd(bin),
		},
	}
}
