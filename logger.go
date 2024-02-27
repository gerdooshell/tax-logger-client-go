package logger

import (
	"context"
	"fmt"
	loggerServer "github.com/gerdooshell/tax-communication/src/logger"
	"github.com/gerdooshell/tax-logger-client-go/internal"
	"github.com/gerdooshell/tax-logger-client-go/internal/environment"
)

func SetUpLogger(ctx context.Context, envStr, configFileAbsPath string) (err error) {
	env, err := environment.GetEnvironmentFromString(envStr)
	if err != nil {
		return err
	}
	if err = internal.ConfigureLoggerByFilePath(ctx, env, configFileAbsPath); err != nil {
		return err
	}
	if _, err = internal.GetClientLoggerInstance(); err != nil {
		err = fmt.Errorf("failed establishing connection to the logging service: %v", err)
	}
	if err == nil {
		fmt.Println("logger is initialized and the logging service is responding")
	}
	return err
}

func Destruct() error {
	return internal.Disconnect()
}

func ErrorWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Error, message, &loggerServer.SaveOriginLogReq{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func ErrorFormat(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}

func ErrorSafe(err error) {
	if err == nil {
		return
	}
	Error(err.Error())
}

func Error(message string) {
	if err := ErrorWithOptions(context.Background(), message, "", ""); err != nil {
		fmt.Println(err)
	}
}

func WarningWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Warning, message, &loggerServer.SaveOriginLogReq{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func WarningFormat(format string, a ...any) {
	Warning(fmt.Sprintf(format, a...))
}

func Warning(message string) {
	if err := WarningWithOptions(context.Background(), message, "", ""); err != nil {
		fmt.Println(err)
	}
}

func InfoWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Info, message, &loggerServer.SaveOriginLogReq{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func InfoFormat(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}

func Info(message string) {
	if err := InfoWithOptions(context.Background(), message, "", ""); err != nil {
		fmt.Println(err)
	}
}

func FatalWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Fatal, message, &loggerServer.SaveOriginLogReq{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func FatalSafe(err error) {
	if err == nil {
		return
	}
	Fatal(err.Error())
}

func FatalFormat(format string, a ...any) {
	Fatal(fmt.Sprintf(format, a...))
}

func Fatal(message string) {
	if err := FatalWithOptions(context.Background(), message, "", ""); err != nil {
		fmt.Println(err)
	}
}

func DebugWithOptions(ctx context.Context, message, stackTrace, processId string) error {
	client, err := internal.GetClientLoggerInstance()
	if err != nil {
		return err
	}
	return client.Log(ctx, internal.Debug, message, &loggerServer.SaveOriginLogReq{
		StackTrace: stackTrace,
		ProcessId:  processId,
	})
}

func DebugFormat(format string, a ...any) {
	Debug(fmt.Sprintf(format, a...))
}

func Debug(message string) {
	if err := DebugWithOptions(context.Background(), message, "", ""); err != nil {
		fmt.Println(err)
	}
}
