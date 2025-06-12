package server

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServerLogger struct {
	logger *logrus.Logger
}

func NewServerLogger(logger *logrus.Logger) *ServerLogger {
	return &ServerLogger{logger: logger}
}

func (s *ServerLogger) LoggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	s.logger.Debugf("Method: %s, Metadata: %s", info.FullMethod, md)

	resp, err := handler(ctx, req)
	s.logger.Debugf("Response: %v, Error: %v\n", resp, err)

	return resp, err
}
