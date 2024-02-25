package logger

import (
	"context"
	loggerServer "github.com/gerdooshell/tax-communication/src/logger"
	"github.com/gerdooshell/tax-logger-client-go/internal"
)

func SetUpLogger(url string) error {
	return internal.SetLoggerConfig(internal.LoggerConfig{Url: url, RegisteredServiceName: "mock"})
}

func Destruct() error {
	return internal.Disconnect()
}

func ErrorWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Error, message, &loggerServer.SaveOriginLogRequest{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func Error(ctx context.Context, message string) error {
	return ErrorWithOptions(ctx, message, "", "")
}

func WarningWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Warning, message, &loggerServer.SaveOriginLogRequest{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func Warning(ctx context.Context, message string) error {
	return WarningWithOptions(ctx, message, "", "")
}

func InfoWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Info, message, &loggerServer.SaveOriginLogRequest{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func Info(ctx context.Context, message string) error {
	return InfoWithOptions(ctx, message, "", "")
}

func FatalWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Fatal, message, &loggerServer.SaveOriginLogRequest{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func Fatal(ctx context.Context, message string) error {
	return FatalWithOptions(ctx, message, "", "")
}

func DebugWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Debug, message, &loggerServer.SaveOriginLogRequest{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func Debug(ctx context.Context, message string) error {
	return DebugWithOptions(ctx, message, "", "")
}
