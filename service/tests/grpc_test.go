package test_grpc

import (
	"context"
	"fmt"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"
	"promos/tests/mock"
	"testing"

	"promos/pkg/postgres"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client promos.PromosClient
var serviceUsers *mock.MockServiceUsers

func init() {

	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Config
	cfg := config.Config{
		Server: server.ServerConfig{
			Network: "tcp",
			Port:    "50123",
		},
		DB: postgres.DBConfig{
			Host:        "localhost",
			Port:        "5432",
			User:        "postgres",
			Password:    "mAz0H1zm",
			Database:    "postgres",
			MaxAtmps:    5,
			DelayAtmpsS: 5,
		},
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

	// Connect to other services
	serviceUsers = mock.NewMockServiceUsers(pool)

	// Creating repository
	repo := repository.NewRepository(serviceUsers)

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

	t.Run("ADDING AND DELETING PROMOCODE", func(t *testing.T) {
		t.Cleanup(func() {
			if _, err := client.DeleteByName(context.TODO(), &promos.PromoName{Name: promo.Name}); err != nil {
				t.Fatal(err)
			}
		})

		if _, err := client.Create(context.TODO(), promo); err != nil {
			t.Fail()
		}

		if _, err := client.Create(context.TODO(), promo); err == nil {
			t.Fail()
		}
	})

	t.Run("USING PROMOCODE", func(t *testing.T) {
		var userId int64

		t.Cleanup(func() {
			if _, err := client.DeleteByName(context.TODO(), &promos.PromoName{Name: promo.Name}); err != nil {
				t.Fail()
			}

			if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
				t.Fail()
			}
		})

		promocode, err := client.GetByName(context.TODO(), &promos.PromoName{Name: promo.Name})
		if err != nil {
			t.Fail()
		}

		userId, err = serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}

		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promocode.PromoCode.Id, UserId: userId}); err != nil {
			t.Fail()
		}
	})
}
