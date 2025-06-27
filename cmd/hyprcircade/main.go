package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	t "time"

	stopCmd "github.com/tfk70/hyprcircade/cmd/hyprcircade/daemon"
	switchCmd "github.com/tfk70/hyprcircade/cmd/hyprcircade/switch"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/tfk70/hyprcircade/internal/time"
	"github.com/tfk70/hyprcircade/pkg/daemon"
	"github.com/tfk70/hyprcircade/pkg/switcher"

	godaemon "github.com/sevlyar/go-daemon"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

const (
	BIN_NAME = "hyprcircade"
	VERSION  = "v0.0.8"
)

func main() {
	cmd := &cli.Command{
		Name:      BIN_NAME,
		Usage:     "Dark/light theme manager for hyprland",
		Version:   VERSION,
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
			&cli.BoolFlag{
				Name:    "foreground",
				Value:   false,
				Usage:   "Run in the foreground mode",
				Local:   true,
				Sources: cli.EnvVars("HYPRCIRCADE_DAEMONIZE"),
			},
		},
		Commands: []*cli.Command{
			switchCmd.SwitchCommand,
			stopCmd.CreateDaemonCmd(BIN_NAME),
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

	runInForeground := cmd.Bool("foreground")

	if !runInForeground {
		dc := &godaemon.Context{
			PidFileName: fmt.Sprintf("/tmp/%s.pid", BIN_NAME),
			PidFilePerm: 0644,
			LogFileName: fmt.Sprintf("/tmp/%s.log", BIN_NAME),
			LogFilePerm: 0640,
		}

		d, err := dc.Reborn()
		if err != nil {
			return err
		}
		if d != nil {
			logger.Info("Daemon started")
			return nil
		}
		defer func() {
			logger.Info("Releasing")
			err := dc.Release()
			if err != nil {
				logger.Errorf("Unable to release file: %v", err)
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			logger.Infof("Received signal %s, shutting down...", sig)
			dc.Release() // remove PID file
			os.Exit(0)
		}()
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
