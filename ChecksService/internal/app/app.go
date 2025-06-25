package app

import (
	"checks/config"
	"checks/internal/repository"
	service "checks/internal/transport/grpc/handlers"
	migrations "checks/pkg/goose"
	"checks/pkg/grpc/conn"
	"checks/pkg/grpc/server"
	"checks/pkg/hasher"
	"checks/pkg/models/checks"
	"context"
	"fmt"

	log "checks/pkg/logger"
	"checks/pkg/postgres"

	"google.golang.org/grpc"
)

func Run() {
	// Setup logger
	logger := log.NewLogger()

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

	// Creating hasher
	hasher := hasher.NewHasher(cfg.Hash)

	// Creating repository
	repo := repository.NewRepository(hasher)

	// Register service promos
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(server.NewServerLogger(logger).LoggingUnaryInterceptor))
	service := service.NewServiceChecks(repo, pool, clientServices)
	checks.RegisterChecksServer(grpcSrv, service)

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
