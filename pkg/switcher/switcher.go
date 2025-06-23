package switcher

import (
	"github.com/tfk70/hyprcircade/internal/commands"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/files"
)

func SwitchToLight(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
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

	return nil
}

func SwitchToDark(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
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

	return nil
}
