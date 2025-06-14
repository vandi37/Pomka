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

	resp, err := handler(ctx, req)
	s.logger.Debugf("METHOD: %s, REQUEST: %s, CONTEXT: %s RESPONSE: %v, ERROR: %v\n", info.FullMethod, req, md, resp, err)

	return resp, err
}
