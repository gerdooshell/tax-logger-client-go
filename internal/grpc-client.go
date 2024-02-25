package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	loggerServer "github.com/gerdooshell/tax-communication/src/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoggerClient struct {
	grpcClient  loggerServer.GRPCLoggerClient
	serverURL   string
	serviceName string
}

var loggerClientInstance *LoggerClient

func GetClientLoggerInstance() (*LoggerClient, error) {
	if loggerClientInstance != nil {
		return loggerClientInstance, nil
	}
	config, err := getLoggerConfig()
	if err != nil {
		return nil, err
	}
	loggerClientInstance = &LoggerClient{
		serverURL:   config.Url,
		serviceName: config.RegisteredServiceName,
	}
	return loggerClientInstance, nil
}

var singletonConnection *grpc.ClientConn
var connectionMu sync.Mutex

func (lc *LoggerClient) generateDataServiceClient() error {
	connectionMu.Lock()
	defer connectionMu.Unlock()
	if lc.grpcClient != nil {
		return nil
	}
	connection, err := grpc.Dial(lc.serverURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		err = fmt.Errorf("connection failed to the logger server")
	}
	lc.grpcClient = loggerServer.NewGRPCLoggerClient(connection)
	singletonConnection = connection
	//if err = connection.Close(); err != nil {
	//	return nil, fmt.Errorf("failed closing connection, error: %v\n", err)
	//}
	return err
}

func Disconnect() error {
	if singletonConnection == nil {
		return errors.New("there is no connection to the logger service")
	}
	return singletonConnection.Close()
}

func (lc *LoggerClient) Log(ctx context.Context, severity Severity, message string, originLog *loggerServer.SaveOriginLogRequest) error {
	if err := lc.generateDataServiceClient(); err != nil {
		return err
	}
	originLog.ServiceName = lc.serviceName
	input := &loggerServer.SaveServiceLogRequest{
		Timestamp: timestamppb.New(time.Now()),
		Severity:  string(severity),
		Message:   message,
		OriginLog: originLog,
	}
	_, err := lc.grpcClient.SaveServiceLog(ctx, input)
	return err
}
