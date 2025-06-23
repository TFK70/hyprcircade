package cron

import (
	"fmt"

	c "github.com/robfig/cron/v3"
)

var cronInstance *c.Cron

func RunEveryMinute(cmd func()) error {
	if cronInstance == nil {
		cronInstance = c.New()
	}

	cronInstance.AddFunc("@every 1m", cmd)

	return nil
}

func Start() error {
	if cronInstance == nil {
		return fmt.Errorf("Failed to start cron: cronInstance was not initialized")
	}

	cronInstance.Start()

	return nil
}
