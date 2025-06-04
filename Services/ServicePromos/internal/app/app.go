package app

import (
	"log"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"

	"google.golang.org/grpc"
)

func Run() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// GRPC server
	grpcSrv := grpc.NewServer()

	// Register promo service
	service := service.NewServicePromos(repository.NewRepository())
	promos.RegisterPromosServer(grpcSrv, service)

	// Run server
	server := server.NewServer(grpcSrv)
	if err := server.Run(cfg.Server); err != nil {
		log.Fatal(err)
	}
}
