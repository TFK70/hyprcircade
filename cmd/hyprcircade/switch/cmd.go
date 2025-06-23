package cmd

import (
	switchDark "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch/dark"
	switchLight "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch/light"
	"github.com/urfave/cli/v3"
)

var (
	SwitchCommand = &cli.Command{
		Name: "switch",
		Commands: []*cli.Command{
			switchDark.SwitchDarkCommand,
			switchLight.SwitchLightCommand,
		},
	}
)
