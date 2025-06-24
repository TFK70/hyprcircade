package daemon

import (
	"github.com/sirupsen/logrus"
	"github.com/tfk70/hyprcircade/internal/config"
	"github.com/tfk70/hyprcircade/internal/cron"
	"github.com/tfk70/hyprcircade/internal/logging"
	"github.com/tfk70/hyprcircade/pkg/switcher"
)

func StartDaemon(darkAt int, lightAt int, cfgFiles []*config.File, cfgCommands []*config.Command, cfgAnchor string) error {
	logger, err := logging.GetNamedLogger("daemon.go")
	if err != nil {
		return err
	}

	switchToDarkFunc := func() {
		switcher.SwitchToDark(cfgFiles, cfgCommands, cfgAnchor)
	}
	switchToLightFunc := func() {
		switcher.SwitchToLight(cfgFiles, cfgCommands, cfgAnchor)
	}

	err = cron.RunEveryNthHour(darkAt, switchToDarkFunc)
	if err != nil {
		return err
	}

	err = cron.RunEveryNthHour(lightAt, switchToLightFunc)
	if err != nil {
		return err
	}

	err = cron.Start()
	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"darkAt":  darkAt,
		"lightAt": lightAt,
	}).Info("Started Daemon")

	return nil
}
