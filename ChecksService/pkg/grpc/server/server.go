package server

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	grpc   *grpc.Server
	logger *logrus.Logger
}

func NewServer(grpc *grpc.Server) *Server {
	return &Server{grpc: grpc}
}

type ServerConfig struct {
	Network string
	Port    string
}

func (s *Server) Run(cfg ServerConfig) error {

	lis, err := net.Listen(cfg.Network, ":"+cfg.Port)
	if err != nil {
		return fmt.Errorf("server: Run: %s", err)
	}

	if err := s.grpc.Serve(lis); err != nil {
		return fmt.Errorf("server: Run: %s", err)
	}

	return nil
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
