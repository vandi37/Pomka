package warns_test

import (
	"context"
	"fmt"
	"testing"
	"warns/config"
	"warns/internal/repository"
	service "warns/internal/transport/grpc/handlers"
	Err "warns/pkg/errors"
	"warns/pkg/grpc/server"
	"warns/pkg/models/users"
	"warns/pkg/models/warns"
	"warns/tests/mock"

	"warns/pkg/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
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
var repo *repository.Repository
var pool *pgxpool.Pool
var cfg config.Config

func init() {

	// Setup logger
	logger = logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	serverLogger := server.NewServerLogger(logger)

	// Load env
	err := godotenv.Load("config.env")
	if err != nil {
		panic("error missing enviroment file. please create config.env in ./service/tests")
	}

	// Cofiguration
	cfg, err = config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Up docker postgres
	dockerpostgres, err = mock.PostgresUp(mock.Config{User: cfg.DB.User, Password: cfg.DB.Password, Name: cfg.DB.Database})
	if err != nil {
		panic(err)
	}

	// Connecting to postgres
	pool, err = postgres.NewPool(context.TODO(), cfg.DB)
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
	repo = repository.NewRepository(logger)

	// Register promo service
	service := service.NewServiceWarns(repo, pool, cfg.Warns, serviceUsers)
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

	t.Run("WARN UNWARN BAN UNBAN", AllTest)
}

func AllTest(t *testing.T) {
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
			err: Err.ErrUserAlreadyBanned,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Send warns before user got banned
			for i := 0; i < cfg.Warns.WarnsBeforeBan; i++ {
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
		logger.Debugf("deleting user: %d", userId)
		if err := serviceUsers.Delete(context.TODO(), userId); err != nil {
			return err
		}
	}

	return nil
}
