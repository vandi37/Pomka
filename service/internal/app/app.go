package app

import (
	"context"
	"log"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"

	"promos/pkg/postgres"

	"google.golang.org/grpc"
)

func Run() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Connecting to postgres
	pool, err := postgres.NewPool(context.TODO(), cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	if err := pool.Ping(context.TODO()); err != nil {
		log.Fatal(err)
	}

	// GRPC server
	grpcSrv := grpc.NewServer()

	// Creating repository
	repo := repository.NewRepository(pool)

	// Register promo service
	service := service.NewServicePromos(repo)
	promos.RegisterPromosServer(grpcSrv, service)

	// Run server
	server := server.NewServer(grpcSrv)
	if err := server.Run(cfg.Server); err != nil {
		log.Fatal(err)
	}
}
