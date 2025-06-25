package dbus

import (
	"fmt"

	d "github.com/godbus/dbus/v5"
	"github.com/tfk70/hyprcircade/internal/logging"
)

var (
	connection *d.Conn
	suspendFns = []func(){}
	awakeFns   = []func(){}
)

func initializeDbus() error {
	if connection != nil {
		return fmt.Errorf("Dbus is already initialized")
	}

	conn, err := d.SystemBus()
	if err != nil {
		return err
	}

	connection = conn

	return nil
}

func startSleepSignalsListener() error {
	if connection == nil {
		return fmt.Errorf("Dbus is not initialized")
	}

	logger, err := logging.GetNamedLogger("dbus.go")
	if err != nil {
		return err
	}

	err = connection.AddMatchSignal(
		d.WithMatchInterface("org.freedesktop.login1.Manager"),
		d.WithMatchMember("PrepareForSleep"),
	)
	if err != nil {
		return err
	}

	c := make(chan *d.Signal, 10)
	connection.Signal(c)

	for signal := range c {
		if len(signal.Body) > 0 {
			suspending := signal.Body[0].(bool)
			if suspending {
				logger.Infof("System suspending, running %d functions", len(suspendFns))
				for _, fn := range suspendFns {
					fn()
				}
			} else {
				logger.Infof("System awaken, running %d functions", len(awakeFns))
				for _, fn := range awakeFns {
					fn()
				}
			}
		}
	}

	return nil
}

func RunOnSuspend(fn func()) error {
	if connection == nil {
		err := initializeDbus()
		if err != nil {
			return err
		}

		logger, err := logging.GetNamedLogger("dbus.go")
		if err != nil {
			return err
		}

		go func() {
			goroutineLogger := logger.WithField("goroutine", "true")

			err = startSleepSignalsListener()
			if err != nil {
				goroutineLogger.Errorf("Error starting listeners %v", err)
			}
		}()
	}

	suspendFns = append(suspendFns, fn)

	return nil
}

func RunOnAwake(fn func()) error {
	if connection == nil {
		err := initializeDbus()
		if err != nil {
			return err
		}

		logger, err := logging.GetNamedLogger("dbus.go")
		if err != nil {
			return err
		}

		go func() {
			goroutineLogger := logger.WithField("goroutine", "true")

			err = startSleepSignalsListener()
			if err != nil {
				goroutineLogger.Errorf("Error starting listeners %v", err)
			}
		}()
	}

	awakeFns = append(awakeFns, fn)

	return nil
}
