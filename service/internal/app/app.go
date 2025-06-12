package app

import (
	"context"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	"promos/internal/transport/grpc/conn"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"

	"promos/pkg/postgres"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Run() {
	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	serverLogger := server.NewServerLogger(logger)

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Creating postgres pool
	pool, err := postgres.NewPool(context.TODO(), cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}
	if err := pool.Ping(context.TODO()); err != nil {
		logger.Fatal(err)
	}

	// Connect to service users
	clientServices, err := conn.NewClientsServices(cfg.Conn)
	if err != nil {
		logger.Fatal(err)
	}

	// Creating repository
	repo := repository.NewRepository(clientServices, logger)

	// Register service promos
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(serverLogger.LoggingUnaryInterceptor))
	service := service.NewServicePromos(repo, pool)
	promos.RegisterPromosServer(grpcSrv, service)

	// Run server
	server := server.NewServer(grpcSrv)
	if err := server.Run(cfg.Server); err != nil {
		logger.Fatal(err)
	}
}
