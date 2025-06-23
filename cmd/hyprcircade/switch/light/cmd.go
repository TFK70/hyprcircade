package cmd

import (
	"context"

	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/pkg/switcher"
	"github.com/urfave/cli/v3"
)

var (
	SwitchLightCommand = &cli.Command{
		Name:     "light",
		Usage:    "switch theme to light",
		Category: "switch",
		Action:   run,
	}
)

func run(context context.Context, cmd *cli.Command) error {
	configPath := cmd.String("config")

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}

	err = switcher.SwitchToLight(cfg.Files, cfg.Commands, cfg.General.Anchor)
	if err != nil {
		return err
	}

	return nil
}
