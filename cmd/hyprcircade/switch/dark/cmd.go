package cmd

import (
	"context"

	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/tfk70/hyprcircade/pkg/switcher"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var (
	SwitchDarkCommand = &cli.Command{
		Name:     "dark",
		Usage:    "switch theme to dark",
		Category: "switch",
		Action:   run,
	}
)

func run(context context.Context, cmd *cli.Command) error {
	logger, err := logging.GetLogger()
	if err != nil {
		return err
	}

	if cmd.Bool("debug") {
		logger.SetLevel(logrus.DebugLevel)
		logger.Debug("Debug logging set")
	}

	configPath := cmd.String("config")

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}

	err = switcher.SwitchToDark(cfg.Files, cfg.Commands, cfg.General.Anchor)
	if err != nil {
		return err
	}

	return nil
}
