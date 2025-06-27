package logger

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func (l *Logger) LoggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	l.WithFields(logrus.Fields{
		"METHOD":   info.FullMethod,
		"REQUEST":  req,
		"RESPONSE": fmt.Sprint(resp, err),
	}).Info("gRPC SERVER")

	return resp, err
}
