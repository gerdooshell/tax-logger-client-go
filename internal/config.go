package internal

import "errors"

type LoggerConfig struct {
	Url                   string
	RegisteredServiceName string
}

var loggerConfig *LoggerConfig

func SetLoggerConfig(conf LoggerConfig) error {
	if loggerConfig != nil {
		return errors.New("logger is already configured")
	}
	loggerConfig = &conf
	return nil
}

func getLoggerConfig() (LoggerConfig, error) {
	if loggerConfig == nil {
		return LoggerConfig{}, errors.New("logger is not configured")
	}
	return *loggerConfig, nil
}
