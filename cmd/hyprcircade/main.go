package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	t "time"

	switchCmd "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/pkg/daemon"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:      "hyprcircade",
		Usage:     "Dark/light theme manager for hyprland",
		Version:   "0.0.1",
		Copyright: fmt.Sprintf("(c) %d TFK70", t.Now().Year()),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Value:   filepath.Join(os.Getenv("HOME"), ".config/hypr/hyprcircade.conf"),
				Usage:   "Path to hyprcircade configuration file",
				Aliases: []string{"c"},
				Sources: cli.EnvVars("HYPRCIRCADE_CONFIGURATION_FILE"),
			},
		},
		Commands: []*cli.Command{
			switchCmd.SwitchCommand,
		},
		Action: run,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func run(context context.Context, cmd *cli.Command) error {
	configPath := cmd.String("config")

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return err
	}

	errch := make(chan error)
	defer close(errch)

	err = daemon.StartDaemon(cfg.General.DarkAt, cfg.General.LightAt, cfg.Files, cfg.Commands, cfg.General.Anchor)

	return <-errch
}
