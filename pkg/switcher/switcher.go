package switcher

import (
	"fmt"

	"github.com/tfk70/hyprcircade/internal/commands"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/files"
	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/tfk70/hyprcircade/internal/time"
)

func SwitchByTod(tod string, cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	logger, err := logging.GetNamedLogger("switcher.go")
	if err != nil {
		return err
	}

	logger.Infof("Setting theme based on time of day: %s", tod)

	if tod == time.LIGHT {
		SwitchToLight(cfgFiles, cfgCommands, cfgAnchor)
		return nil
	}

	if tod == time.DARK {
		SwitchToDark(cfgFiles, cfgCommands, cfgAnchor)
		return nil
	}

	return fmt.Errorf("Undefined time of day: %s", tod)
}

func SwitchToLight(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	logger, err := logging.GetNamedLogger("switcher.go")
	if err != nil {
		return err
	}

	lightLogger := logger.WithField("theme", "light")
	lightLogger.Info("Switching theme")

	for _, file := range cfgFiles {
		var anchor string

		if !file.IgnoreAnchor {
			anchor = cfgAnchor
		}

		err := files.ReplaceInFile(file.Path, file.NightValue, file.DayValue, anchor)
		if err != nil {
			return err
		}
	}

	for _, cmd := range cfgCommands {
		if cmd.DayExec != "" {
			err := commands.ExecuteCommand(cmd.DayExec)
			if err != nil {
				return err
			}
		}
	}

	lightLogger.Info("Theme switched successfully")

	return nil
}

func SwitchToDark(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	logger, err := logging.GetNamedLogger("switcher.go")
	if err != nil {
		return err
	}

	darkLogger := logger.WithField("theme", "dark")
	darkLogger.Info("Switching theme")

	for _, file := range cfgFiles {
		var anchor string

		if !file.IgnoreAnchor {
			anchor = cfgAnchor
		}

		err := files.ReplaceInFile(file.Path, file.DayValue, file.NightValue, anchor)
		if err != nil {
			return err
		}
	}

	for _, cmd := range cfgCommands {
		if cmd.NightExec != "" {
			err := commands.ExecuteCommand(cmd.NightExec)
			if err != nil {
				return err
			}
		}
	}

	darkLogger.Info("Theme switched successfully")

	return nil
}
