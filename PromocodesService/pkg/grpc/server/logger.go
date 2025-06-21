package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ServerLogger struct {
	logger *logrus.Logger
}

func NewServerLogger(logger *logrus.Logger) *ServerLogger {
	return &ServerLogger{logger: logger}
}

func (s *ServerLogger) LoggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	s.logger.WithFields(logrus.Fields{
		"METHOD":   info.FullMethod,
		"REQUEST":  req,
		"RESPONSE": resp,
		"ERROR":    err,
	}).Info("gRPC SERVER")

	return resp, err
}
