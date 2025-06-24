package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetupLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}

	logger = logrus.New()

	return logger
}

func GetLogger() (*logrus.Logger, error) {
	if logger == nil {
		return nil, fmt.Errorf("Logger not found")
	}

	return logger, nil
}

func GetNamedLogger(name string) (*logrus.Entry, error) {
	if name == "" {
		return nil, fmt.Errorf("Name for logger was not provided")
	}

	if logger == nil {
		return nil, fmt.Errorf("Logger not found")
	}

	namedLogger := logger.WithField("logger", name)

	return namedLogger, nil
}
