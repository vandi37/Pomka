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
var repo *repository.Repository

func init() {

	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	serverLogger := server.NewServerLogger(logger)

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
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(serverLogger.LoggingUnaryInterceptor))

	// Connect to other services
	serviceUsers = mock.NewMockServiceUsers(pool)

	// Creating repository
	repo = repository.NewRepository(serviceUsers, logger)

	// Register promo service
	service := service.NewServicePromos(repo, pool)
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
		var promoId, userId int64

		t.Cleanup(func() {
			if _, err := client.DeleteById(context.TODO(), &promos.PromoId{Id: promoId}); err != nil {
				t.Fatal(err)
			}

			if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
				t.Fatal(err)
			}
		})

		userId, err := serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}
		promo.Creator = userId

		outCreate, err := client.Create(context.TODO(), promo)
		if err != nil {
			t.Fail()
		}
		promoId = outCreate.PromoCode.Id

		if _, err := client.Create(context.TODO(), promo); err == nil {
			t.Fail()
		}
	})

	t.Run("USING PROMOCODE", func(t *testing.T) {
		var userId, promoId int64

		t.Cleanup(func() {

			if err := serviceUsers.ClearHistory(context.TODO(), userId, promoId); err != nil {
				t.Fail()
			}

			if _, err := client.DeleteById(context.TODO(), &promos.PromoId{Id: promoId}); err != nil {
				t.Fail()
			}

			if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
				t.Fail()
			}

		})

		userId, err := serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}
		promo.Creator = userId

		promoFailure, err := client.Create(context.TODO(), promo)
		if err != nil {
			t.Fail()
		}
		promoId = promoFailure.PromoCode.Id

		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId}); err != nil {
			t.Fail()
		}
	})
}
