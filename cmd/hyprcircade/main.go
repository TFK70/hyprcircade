package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	t "time"

	switchCmd "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/tfk70/hyprcircade/internal/time"
	"github.com/tfk70/hyprcircade/pkg/daemon"
	"github.com/tfk70/hyprcircade/pkg/switcher"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "hyprcircade",
		Usage:     "Dark/light theme manager for hyprland",
		Version:   "v0.0.7",
		Copyright: fmt.Sprintf("(c) %d TFK70", t.Now().Year()),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   filepath.Join(os.Getenv("HOME"), ".config/hypr/hyprcircade.conf"),
				Usage:   "Path to hyprcircade configuration file",
				Aliases: []string{"c"},
				Sources: cli.EnvVars("HYPRCIRCADE_CONFIGURATION_FILE"),
			},
			&cli.BoolFlag{
				Name:    "debug",
				Value:   false,
				Usage:   "Enable debug logging",
				Sources: cli.EnvVars("HYPRCIRCADE_DEBUG"),
			},
			&cli.BoolFlag{
				Name:    "apply-on-start",
				Value:   true,
				Usage:   "Apply theme based on time of day on daemon startup",
				Local:   true,
				Sources: cli.EnvVars("HYPRCIRCADE_APPLY_ON_START"),
			},
		},
		Commands: []*cli.Command{
			switchCmd.SwitchCommand,
		},
		Action: run,
	}

	logger := logging.SetupLogger()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		if logging.IsNullified {
			fmt.Println(fmt.Errorf("Error during execution: %v", err))
		} else {
			logger.Errorf("Error during execution: %v", err)
		}

		os.Exit(1)
	}
}

func run(context context.Context, cmd *cli.Command) error {
	rootLogger, err := logging.GetLogger()
	if err != nil {
		return err
	}

	if cmd.Bool("debug") {
		rootLogger.SetLevel(logrus.DebugLevel)
	}

	logger, err := logging.GetNamedLogger("main.go")
	if err != nil {
		return err
	}

	configPath := cmd.String("config")

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}

	errch := make(chan error)
	defer close(errch)

	err = daemon.StartDaemon(cfg.General.DarkAt, cfg.General.LightAt, cfg.Files, cfg.Commands, cfg.General.Anchor)
	if err != nil {
		return err
	}

	if cmd.Bool("apply-on-start") {
		logger.Info("apply-on-start flag set, applying theme based on time of day")
		tod, err := time.GetCurrentTimeOfTheDay(cfg.General.DarkAt, cfg.General.LightAt)
		if err != nil {
			return err
		}
		switcher.SwitchByTod(tod, cfg.Files, cfg.Commands, cfg.General.Anchor)
	}

	return <-errch
}
