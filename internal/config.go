package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gerdooshell/tax-logger-client-go/internal/environment"
	"github.com/gerdooshell/tax-logger-client-go/internal/secret-service/azure"
)

type LoggerConfig struct {
	Url                   string
	RegisteredServiceName string
	APIKey                string
}

var loggerConfig *LoggerConfig

func setLoggerConfig(conf LoggerConfig) error {
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

type loggerConfigModel struct {
	ContainsSecretKeys bool   `json:"ContainsSecretKeys"`
	VaultURL           string `json:"VaultURL"`
	LoggerUrl          string `json:"LoggerUrl"`
	ServiceName        string `json:"ServiceName"`
	APIKey             string `json:"APIKey"`
	Port               string `json:"Port"`
}

var isLoggerConfigured = false

func ConfigureLoggerByFilePath(ctx context.Context, env environment.Environment, absFilePath string) error {
	if isLoggerConfigured {
		return errors.New("logger is already configured")
	}
	config, err := getLoggerConfigModel(ctx, absFilePath, env)
	if err != nil {
		return err
	}
	loggerHost := fmt.Sprintf("%s:%s", config.LoggerUrl, config.Port)
	if err = setLoggerConfig(LoggerConfig{
		Url:                   loggerHost,
		RegisteredServiceName: config.ServiceName,
		APIKey:                config.APIKey}); err != nil {
		return err
	}
	isLoggerConfigured = true
	return nil
}

func getLoggerConfigModel(ctx context.Context, absFilePath string, env environment.Environment) (*loggerConfigModel, error) {
	data, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, err
	}
	var confMap map[environment.Environment]loggerConfigModel
	if err = json.Unmarshal(data, &confMap); err != nil {
		return nil, err
	}
	conf, ok := confMap[env]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no config found for environment %v", env))
	}
	if conf.ContainsSecretKeys {
		conf.VaultURL = strings.Trim(conf.VaultURL, " ")
		if conf.VaultURL == "" {
			return nil, errors.New("invalid vault url")
		}
		if conf, err = setVaultSecrets(ctx, conf, env); err != nil {
			return nil, err
		}
	}
	return &conf, nil
}

func setVaultSecrets(ctx context.Context, conf loggerConfigModel, env environment.Environment) (loggerConfigModel, error) {
	timeout := time.Second * 15
	azVault := azure.NewSecretService(conf.VaultURL, env)
	loggerUrlChan := azVault.GetSecretValue(ctx, conf.LoggerUrl)
	apiKeyChan := azVault.GetSecretValue(ctx, conf.APIKey)
	ServiceNameChan := azVault.GetSecretValue(ctx, conf.ServiceName)
	portChan := azVault.GetSecretValue(ctx, conf.Port)
	select {
	case secretOut := <-loggerUrlChan:
		if secretOut.Err != nil {
			return conf, secretOut.Err
		}
		conf.LoggerUrl = secretOut.Value
	case <-time.After(timeout):
		return conf, errors.New("fetching logger url secret timed out")
	}
	select {
	case secretOut := <-apiKeyChan:
		if secretOut.Err != nil {
			return conf, secretOut.Err
		}
		conf.APIKey = secretOut.Value
	case <-time.After(timeout):
		return conf, errors.New("fetching logger api key timed out")
	}
	select {
	case secretOut := <-ServiceNameChan:
		if secretOut.Err != nil {
			return conf, secretOut.Err
		}
		conf.ServiceName = secretOut.Value
	case <-time.After(timeout):
		return conf, errors.New("fetching logger service name timed out")
	}
	select {
	case secretOut := <-portChan:
		if secretOut.Err != nil {
			return conf, secretOut.Err
		}
		conf.Port = secretOut.Value
	case <-time.After(timeout):
		return conf, errors.New("fetching logger port secret timed out")
	}
	return conf, nil
}
