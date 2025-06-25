package cheks_test

import (
	"checks/config"
	"checks/internal/repository"
	service "checks/internal/transport/grpc/handlers"
	migrations "checks/pkg/goose"
	"checks/pkg/grpc/server"
	"checks/pkg/hasher"
	"checks/pkg/models/checks"
	"checks/pkg/models/common"
	"checks/pkg/models/users"
	"checks/tests/mock"
	"context"
	"fmt"
	"testing"

	"checks/pkg/postgres"

	log "checks/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client checks.ChecksClient
var serviceUsers *mock.MockServiceUsers
var dockerpostgres *mock.DockerPool
var logger *logrus.Logger
var repo *repository.Repository
var pool *pgxpool.Pool
var cfg config.Config

func TestMain(m *testing.M) {

	defer func() {

		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := dockerpostgres.PostgresDown(); err != nil {
			logger.Fatal(err)
		}

	}()

	// Setup logger
	logger = log.NewLogger()

	// Load enviroment
	if err := godotenv.Load("config.env"); err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	logger.WithField("MSG", "Succecs loading enviroment file").Debug("SETUP APP")

	// Cofiguration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	logger.WithField("MSG", "Succecs loading configuration for app").Debug("SETUP APP")

	// Up docker postgres
	dockerpostgres, err = mock.PostgresUp(mock.Config{User: cfg.DB.User, Password: cfg.DB.Password, Name: cfg.DB.Database})
	if err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	logger.WithField("MSG", "Succecs up postgres docker container").Debug("SETUP APP")

	// Connecting to postgres
	pool, err = postgres.NewPool(context.TODO(), cfg.DB)
	if err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	logger.WithField("MSG", "Succecs connect to postgres").Debug("SETUP APP")

	// Run migrations
	if err := migrations.Up(context.TODO(), pool); err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	logger.WithField("MSG", "Succecs run migrations").Debug("SETUP APP")

	// gRPC server
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(server.NewServerLogger(logger).LoggingUnaryInterceptor))

	// Creating mock service users
	serviceUsers = mock.NewMockServiceUsers(pool)

	// Creating hasher
	hasher := hasher.NewHasher(cfg.Hash)

	// Creating repository
	repo = repository.NewRepository(hasher)

	// Register promo service
	service := service.NewServiceChecks(repo, pool, serviceUsers)
	checks.RegisterChecksServer(grpcSrv, service)

	// Run server
	srv = server.NewServer(grpcSrv)
	go func() {
		if err := srv.Run(cfg.Server); err != nil {
			logger.WithField("ERROR", err).Panic("SETUP APP")
		}
	}()
	logger.WithField("MSG", fmt.Sprintf("Running server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("SETUP APP")

	// Connection to server
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", cfg.Server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.WithField("ERROR", err).Panic("SETUP APP")
	}
	client = checks.NewChecksClient(conn)
	logger.WithField("MSG", fmt.Sprintf("Succecs connection to server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("SETUP APP")

	m.Run()
}

func TestAll(t *testing.T) {
	var creatorId int64

	t.Cleanup(func() {
		usersIds := []int64{creatorId}

		if err := clearUsers(usersIds); err != nil {
			t.Fatal()
		}
	})

	creatorId, err := serviceUsers.Create(context.TODO(), 1)
	if err != nil {
		t.Fatal()
	}

	var tests = []struct {
		name string
		in   *checks.CheckCreate
	}{
		{
			name: "common",
			in: &checks.CheckCreate{
				Creator:  creatorId,
				Currency: common.Currency_Credits,
				Amount:   999,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create check
			check, err := client.Create(context.TODO(), tt.in)
			if err != nil {
				t.Fail()
			}

			// Use check
			if _, err := client.Use(context.TODO(), &checks.CheckUse{Key: check.Check.Key, UserId: creatorId}); err != nil {
				t.Fail()
			}

			// Get checks created by user
			allChecksFailure, err := client.GetUserChecks(context.TODO(), &users.Id{Id: creatorId})
			if err != nil {
				t.Fail()
			}

			// Proof check deleting after using
			if allChecksFailure.AllChecks.Checks != nil {
				t.Fail()
			}
		})
	}
}

func clearUsers(userIds []int64) error {
	for _, userId := range userIds {

		if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
			return err
		}
	}

	return nil
}
