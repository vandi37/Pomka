package cheks_test

import (
	"checks/tests/mock"
	"config"
	"context"
	"fmt"
	"market/internal/repository"
	service "market/internal/transport/grpc/handlers"
	"migrations"
	"protobuf/checks"
	"protobuf/market"
	"server"
	"testing"

	"postgres"

	log "logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client checks.ChecksClient
var serviceUsers *mock.MockServiceUsers
var dockerpostgres *mock.DockerPool
var repo *repository.Repository
var pool *pgxpool.Pool

func TestMain(m *testing.M) {

	defer func() {

		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := dockerpostgres.PostgresDown(); err != nil {
			panic(err)
		}

	}()

	// Setup logger
	logger := log.NewLogger()

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
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(logger.LoggingUnaryInterceptor))

	// Creating mock service users
	serviceUsers = mock.NewMockServiceUsers(pool)

	// Creating repository
	repo = repository.NewRepository()

	// Register promo service
	service := service.NewServiceMarket(repo, pool, serviceUsers)
	market.RegisterMarketServer(grpcSrv, service)

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
