package daemon

import (
	"fmt"

	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/cron"
	"github.com/tfk70/hyprcircade/internal/time"
	"github.com/tfk70/hyprcircade/pkg/switcher"
)

func StartDaemon(darkAt int, lightAt int, cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	reconcileFunc := func() {
		tod, err := time.GetCurrentTimeOfTheDay(darkAt, lightAt)
		if err != nil {
			fmt.Println(fmt.Errorf("Error getting current time of day: %v", err))
		}

		if tod == time.DARK {
			switcher.SwitchToDark(cfgFiles, cfgCommands, cfgAnchor)
		} else if tod == time.LIGHT {
			switcher.SwitchToLight(cfgFiles, cfgCommands, cfgAnchor)
		} else {
			fmt.Println(fmt.Errorf("Undefined time of day value: %s", tod))
		}
	}

	err := cron.RunEveryMinute(reconcileFunc)
	if err != nil {
		return err
	}
	err = cron.Start()
	if err != nil {
		return err
	}

	return nil
}
