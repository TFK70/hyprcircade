package cmd

import (
	switchDark "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch/dark"
	switchLight "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch/light"
	"github.com/urfave/cli/v3"
)

var (
	SwitchCommand = &cli.Command{
		Name: "switch",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "disable-tui",
				Value:   false,
				Usage:   "Use standard logs instead tui",
				Sources: cli.EnvVars("HYPRCIRCADE_DISABLE_TUI"),
			},
		},
		Commands: []*cli.Command{
			switchDark.SwitchDarkCommand,
			switchLight.SwitchLightCommand,
		},
	}
)
