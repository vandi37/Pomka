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

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client promos.PromosClient

func init() {
	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Configuration
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
	repo := repository.NewRepository()

	// Register promo service
	service := service.NewServicePromos(repo, pool, logger)
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

func Test(t *testing.T) {
	t.Cleanup(func() {
		srv.Stop()
	})

	t.Run("Test Create|Use|Delete", func(t *testing.T) {
		if _, err := client.Create(context.TODO(), promo); err != nil {
			t.Fatal(err)
		}

		if _, err := client.Delete(context.TODO(), &promos.PromoName{Name: promo.Name}); err != nil {
			t.Fatal(err)
		}
	})

}
