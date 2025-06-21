package app

import (
	"context"
	"fmt"
	"promos/config"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	"promos/pkg/grpc/conn"
	"promos/pkg/grpc/server"
	"promos/pkg/models/promos"

	"promos/pkg/postgres"

	"promos/pkg/logger"

	"google.golang.org/grpc"
)

func Run() {
	// Setup logger
	logger := logger.NewLogger()

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithField("ERROR", err).Fatal("SETUP APP")
	}
	logger.WithField("MSG", "Succecs loading configuration for app").Debug("SETUP APP")

	// Creating postgres pool
	pool, err := postgres.NewPool(context.TODO(), cfg.DB)
	if err != nil {
		logger.WithField("ERROR", err).Fatal("SETUP APP")
	}
	if err := pool.Ping(context.TODO()); err != nil {
		logger.WithField("ERROR", err).Fatal("SETUP APP")
	}
	logger.WithField("MSG", "Succecs connect to postgres").Debug("SETUP APP")

	// Connect to service users
	clientServices, err := conn.NewClientsServices(cfg.Conn)
	if err != nil {
		logger.WithField("ERROR", err).Fatal("SETUP APP")
	}
	logger.WithField("MSG", fmt.Sprintf("Succecs connect to gRPC server (service Users) on %s:%s", cfg.Conn.CfgSrvUsers.Host, cfg.Conn.CfgSrvUsers.Port)).Debug("SETUP APP")

	// Creating repository
	repo := repository.NewRepository()

	// Register service promos
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(server.NewServerLogger(logger).LoggingUnaryInterceptor))
	service := service.NewServicePromos(repo, pool, clientServices)
	promos.RegisterPromosServer(grpcSrv, service)

	// Run server
	logger.WithField("MSG", fmt.Sprintf("Running server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("SETUP APP")

	server := server.NewServer(grpcSrv)
	if err := server.Run(cfg.Server); err != nil {
		logger.WithField("ERROR", err).Fatal("SETUP APP")
	}
}
