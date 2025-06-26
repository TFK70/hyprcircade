package switcher

import (
	"fmt"

	"github.com/tfk70/hyprcircade/internal/commands"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/files"
	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/tfk70/hyprcircade/internal/time"
	"github.com/tfk70/hyprcircade/internal/tui"
)

var (
	withTui    = false
	switchView tui.SwitchView
)

func WithTui() {
	withTui = true
}

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

func SwitchToLightWithTui(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	if !withTui {
		return fmt.Errorf("withTui flag is not set")
	}

	filesSteps := []tui.SwitchModelStepDto{}
	for _, file := range cfgFiles {
		filesSteps = append(filesSteps, tui.SwitchModelStepDto{
			Name:         file.Path,
			PendingMsg:   fmt.Sprintf("Processing file %s", file.Path),
			CompletedMsg: fmt.Sprintf("Processed file %s", file.Path),
		})
	}
	filesStage := tui.SwitchModelStageDto{
		CompletedMsg: "Files processed",
		Steps:        filesSteps,
	}

	commandsSteps := []tui.SwitchModelStepDto{}

	for _, cmd := range cfgCommands {
		if cmd.DayExec != "" {
			commandsSteps = append(commandsSteps, tui.SwitchModelStepDto{
				Name:         cmd.DayExec,
				PendingMsg:   fmt.Sprintf("Executing command %s", cmd.DayExec),
				CompletedMsg: fmt.Sprintf("Executed command %s", cmd.DayExec),
			})
		}
	}
	commandsStage := tui.SwitchModelStageDto{
		CompletedMsg: "Commands executed",
		Steps:        commandsSteps,
	}

	modelStages := []tui.SwitchModelStageDto{filesStage, commandsStage}

	switchView = *tui.CreateSwitchView(modelStages)

	launchedCh := make(chan bool, 1)

	go func() {
		<-launchedCh
		SwitchToLight(cfgFiles, cfgCommands, cfgAnchor)
	}()

	err := switchView.Launch(launchedCh)
	if err != nil {
		return err
	}

	return nil
}

func SwitchToLight(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	logger, err := logging.GetNamedLogger("switcher.go")
	if err != nil {
		return err
	}

	lightLogger := logger.WithField("theme", "light")
	lightLogger.Info("Switching theme")

	for _, file := range cfgFiles {
		if withTui {
			switchView.Pending(file.Path)
		}

		var anchor string

		if !file.IgnoreAnchor {
			anchor = cfgAnchor
		}

		err := files.ReplaceInFile(file.Path, file.NightValue, file.DayValue, anchor)
		if err != nil {
			return err
		}

		if withTui {
			switchView.Proceed(file.Path)
		}
	}

	for _, cmd := range cfgCommands {
		if cmd.DayExec != "" {
			if withTui {
				switchView.Pending(cmd.DayExec)
			}

			err := commands.ExecuteCommand(cmd.DayExec)
			if err != nil {
				return err
			}

			if withTui {
				switchView.Proceed(cmd.DayExec)
			}
		}
	}

	lightLogger.Info("Theme switched successfully")

	return nil
}

func SwitchToDarkWithTui(cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	if !withTui {
		return fmt.Errorf("withTui flag is not set")
	}

	filesSteps := []tui.SwitchModelStepDto{}
	for _, file := range cfgFiles {
		filesSteps = append(filesSteps, tui.SwitchModelStepDto{
			Name:         file.Path,
			PendingMsg:   fmt.Sprintf("Processing file %s", file.Path),
			CompletedMsg: fmt.Sprintf("Processed file %s", file.Path),
		})
	}
	filesStage := tui.SwitchModelStageDto{
		CompletedMsg: "Files processed",
		Steps:        filesSteps,
	}

	commandsSteps := []tui.SwitchModelStepDto{}

	for _, cmd := range cfgCommands {
		if cmd.NightExec != "" {
			commandsSteps = append(commandsSteps, tui.SwitchModelStepDto{
				Name:         cmd.NightExec,
				PendingMsg:   fmt.Sprintf("Executing command %s", cmd.NightExec),
				CompletedMsg: fmt.Sprintf("Executed command %s", cmd.NightExec),
			})
		}
	}
	commandsStage := tui.SwitchModelStageDto{
		CompletedMsg: "Commands executed",
		Steps:        commandsSteps,
	}

	modelStages := []tui.SwitchModelStageDto{filesStage, commandsStage}

	switchView = *tui.CreateSwitchView(modelStages)

	launchedCh := make(chan bool, 1)

	go func() {
		<-launchedCh
		SwitchToDark(cfgFiles, cfgCommands, cfgAnchor)
	}()

	err := switchView.Launch(launchedCh)
	if err != nil {
		return err
	}

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
		if withTui {
			switchView.Pending(file.Path)
		}

		var anchor string

		if !file.IgnoreAnchor {
			anchor = cfgAnchor
		}

		err := files.ReplaceInFile(file.Path, file.DayValue, file.NightValue, anchor)
		if err != nil {
			return err
		}

		if withTui {
			switchView.Proceed(file.Path)
		}
	}

	for _, cmd := range cfgCommands {
		if cmd.NightExec != "" {
			if withTui {
				switchView.Pending(cmd.NightExec)
			}

			err := commands.ExecuteCommand(cmd.NightExec)
			if err != nil {
				return err
			}

			if withTui {
				switchView.Proceed(cmd.NightExec)
			}
		}
	}

	darkLogger.Info("Theme switched successfully")

	return nil
}
