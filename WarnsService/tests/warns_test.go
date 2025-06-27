package warns_test

import (
	"config"
	"context"
	e "errorspomka"
	"fmt"
	"protobuf/users"
	"protobuf/warns"
	"server"
	"testing"
	"warns/internal/repository"
	service "warns/internal/transport/grpc/handlers"
	"warns/tests/mock"

	"postgres"

	log "logger"

	"migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client warns.WarnsClient
var serviceUsers *mock.MockServiceUsers
var dockerpostgres *mock.DockerPool
var repo *repository.Repository
var pool *pgxpool.Pool
var cfg config.Config

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
	service := service.NewServiceWarns(repo, pool, service.Config{WarnsBeforeBan: cfg.Storage.WarnsBeforeBan}, serviceUsers)
	warns.RegisterWarnsServer(grpcSrv, service)

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
	client = warns.NewWarnsClient(conn)
	logger.WithField("MSG", fmt.Sprintf("Succecs connection to server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("SETUP APP")

	m.Run()
}

func TestAll(t *testing.T) {
	var moderId, userId int64

	t.Cleanup(func() {
		userIds := []int64{moderId, userId}

		if err := clearWarnsBans(userIds); err != nil {
			t.Fatal(err)
		}

		if err := clearUsers(userIds); err != nil {
			t.Fatal(err)
		}
	})

	// Create moderator
	moderId, err := serviceUsers.Create(context.TODO(), 2)
	if err != nil {
		t.Fatal(err)
	}

	// Create bad boy
	userId, err = serviceUsers.Create(context.TODO(), 1)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name string
		in   *warns.ModerUserReason
		err  error
	}{
		{
			name: "common",
			in: &warns.ModerUserReason{
				ModerId: moderId,
				UserId:  userId,
				Reason:  nil,
			},
			err: e.ErrUserAlreadyBanned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Send warns before user got banned
			for i := 0; i < cfg.Storage.WarnsBeforeBan; i++ {
				if _, err := client.Warn(context.TODO(), tt.in); err != nil {
					t.Fail()
				}
			}

			// Try ban user, want error, user already banned
			if _, err := client.Ban(context.TODO(), tt.in); !(fmt.Errorf("rpc error: code = Unknown desc = %s", tt.err) != err) {
				t.Fail()
			}

			// Unban user
			if _, err := client.Unban(context.TODO(), tt.in); err != nil {
				t.Fail()
			}

			// Send warn
			if _, err := client.Warn(context.TODO(), tt.in); err != nil {
				t.Fail()
			}

			// Unwarn
			if _, err := client.LastUnWarn(context.TODO(), tt.in); err != nil {
				t.Fail()
			}

			// Check count of warns
			countOfWarns, err := client.GetCountOfActiveWarns(context.TODO(), &users.Id{Id: tt.in.UserId})
			if err != nil || countOfWarns.CountWarns != 0 {
				t.Fail()
			}
		})
	}
}

func clearWarnsBans(userIds []int64) error {
	for _, userId := range userIds {
		if err := repo.DeleteHistoryWarns(context.TODO(), pool, &users.Id{Id: userId}); err != nil {
			return err
		}

		if err := repo.DeleteHistoryBans(context.TODO(), pool, &users.Id{Id: userId}); err != nil {
			return err
		}
	}

	return nil
}

func clearUsers(userIds []int64) error {
	for _, userId := range userIds {

		if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
			return err
		}
	}

	return nil
}
