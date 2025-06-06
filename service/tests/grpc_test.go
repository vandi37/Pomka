package test_grpc

import (
	"context"
	"fmt"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"
	"testing"

	"promos/pkg/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client promos.PromosClient

func init() {
	// Config
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Connecting to postgres
	pool, err := postgres.NewPool(context.TODO(), cfg.DB)
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(context.TODO()); err != nil {
		panic(err)
	}
	// GRPC server
	grpcSrv := grpc.NewServer()

	// Creating repository
	repo := repository.NewRepository(pool)

	// Register promo service
	service := service.NewServicePromos(repo)
	promos.RegisterPromosServer(grpcSrv, service)

	// Run server
	srv = server.NewServer(grpcSrv)
	go func() {
		if err := srv.Run(cfg.Server); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", cfg.Server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = promos.NewPromosClient(conn)
}

func TestAllTests(t *testing.T) {
	t.Cleanup(func() {
		srv.Stop()
	})

	t.Run("TestCreate", create)
	t.Run("TestDelete", delete)
	t.Run("TestUse", use)
}

func create(t *testing.T) {
	out, err := client.Create(context.TODO(), &promos.CreatePromoIn{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("create: ", out)
}

func delete(t *testing.T) {
	out, err := client.Delete(context.TODO(), &promos.PromoName{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("delete: ", out)
}

func use(t *testing.T) {
	out, err := client.Use(context.TODO(), &promos.PromoName{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("use: ", out)
}
