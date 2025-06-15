package promos_test

import (
	"context"
	"fmt"
	"testing"
	"warns/config"
	"warns/internal/repository"
	service "warns/internal/transport/grpc/handlers"
	"warns/pkg/grpc/server"
	"warns/pkg/models/warns"
	"warns/tests/mock"

	"warns/pkg/postgres"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client warns.WarnsClient
var serviceUsers *mock.MockServiceUsers
var dockerpostgres *mock.DockerPool
var logger *logrus.Logger

// ДОПИСАТЬ ТЕСТЫ

func init() {

	// Setup logger
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	serverLogger := server.NewServerLogger(logger)

	// Load env
	if err := godotenv.Load("config.env"); err != nil {
		panic("error missing enviroment file. please create config.env in ./service/tests")
	}

	// Cofiguration
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Up docker postgres
	dockerpostgres, err = mock.PostgresUp(mock.Config{User: cfg.DB.User, Password: cfg.DB.Password, Name: cfg.DB.Database})
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

	// gRPC server
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(serverLogger.LoggingUnaryInterceptor))

	// Creating mock service users
	serviceUsers = mock.NewMockServiceUsers(pool)

	// Creating repository
	repo := repository.NewRepository(serviceUsers, logger)

	// Register promo service
	service := service.NewServiceWarns(repo, pool)
	warns.RegisterWarnsServer(grpcSrv, service)

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
	client = warns.NewWarnsClient(conn)
}

func TestMain(t *testing.T) {
	t.Cleanup(func() {

		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := dockerpostgres.PostgresDown(); err != nil {
			t.Fatal(err)
		}

	})
}

func clearUsers(userIds []int64) error {
	for _, userId := range userIds {
		logger.Debugf("deleting user: %d", userId)
		if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
			return err
		}
	}

	return nil
}
