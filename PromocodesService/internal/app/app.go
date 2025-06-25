package app

import (
	"context"
	"fmt"
	"promos/config"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	migrations "promos/pkg/goose"
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
	logger.WithField("MSG", fmt.Sprintf("Succecs connect to postgres://%s:%s@%s:%s/%s",
		cfg.DB.User, "<PASSWORD>", cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)).Debug("SETUP APP")

	// Run migrations
	if err := migrations.Up(context.TODO(), pool); err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	logger.WithField("MSG", "Succecs run migrations").Debug("SETUP APP")

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

	// defer all stoping
	defer func() {
		pool.Close()
		logger.WithField("MSG", fmt.Sprintf("Closing connect to postgres://%s:%s@%s:%s/%s",
			cfg.DB.User, "<PASSWORD>", cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)).Debug("CLOSING APP")

		if err := clientServices.ClientConn.Close(); err != nil {
			logger.WithField("ERROR", err).Fatal("CLOSING APP")
		}
		logger.WithField("MSG", fmt.Sprintf("Closing connect to gRPC server (service Users) on %s:%s",
			cfg.Conn.CfgSrvUsers.Host, cfg.Conn.CfgSrvUsers.Port)).Debug("CLOSING APP")

		server.Stop()
		logger.WithField("MSG", fmt.Sprintf("Closing server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("CLOSING APP")
	}()

	if err := server.Run(cfg.Server); err != nil {
		logger.WithField("ERROR", err).Fatal("SETUP APP")
	}
}
