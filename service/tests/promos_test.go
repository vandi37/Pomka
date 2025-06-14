package promos_test

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
	"time"

	"promos/pkg/postgres"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var srv *server.Server
var client promos.PromosClient
var serviceUsers *mock.MockServiceUsers
var repo *repository.Repository
var dockerpostgres *mock.DockerPool
var logger *logrus.Logger
var pool *pgxpool.Pool

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

func TestMain(t *testing.T) {
	t.Cleanup(func() {

		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := dockerpostgres.PostgresDown(); err != nil {
			t.Fatal(err)
		}

	})

	t.Run("CREATE", Create)
	t.Run("Use", Use)
}

func Create(t *testing.T) {
	var creator, user int64
	var promoIds []int64

	t.Cleanup(
		func() {
			// Delete testing data from table Promos
			if err := clearPromos(promoIds); err != nil {
				t.Fatal(err)
			}

			// Delete testing data from table Users
			if err := clearUsers([]int64{creator, user}); err != nil {
				t.Fatal(err)
			}
		},
	)

	// Creating creator user
	creator, err := serviceUsers.Create(context.TODO(), 3)
	if err != nil {
		t.Fatal(err)
	}

	// Creating common user
	user, err = serviceUsers.Create(context.TODO(), 1)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name string
		in   *promos.CreatePromo
		err  bool
	}{
		{
			name: "common",
			in: &promos.CreatePromo{
				Name:    uuid.NewString(),
				Uses:    -1,
				ExpAt:   timestamppb.New(time.Now().Add(time.Hour * 12)),
				Creator: creator,
			},
			err: false,
		},
		{
			name: "bad arg uses",
			in: &promos.CreatePromo{
				Name:    uuid.NewString(),
				Uses:    0,
				ExpAt:   timestamppb.New(time.Now().Add(time.Hour * 12)),
				Creator: creator,
			},
			err: true,
		},
		{
			name: "bad arg currency",
			in: &promos.CreatePromo{
				Name:     uuid.NewString(),
				Uses:     -1,
				Currency: 10,
				ExpAt:    timestamppb.New(time.Now().Add(time.Hour * 12)),
				Creator:  creator,
			},
			err: true,
		},
		{
			name: "bad arg amount",
			in: &promos.CreatePromo{
				Name:    uuid.NewString(),
				Uses:    -1,
				Amount:  -10,
				ExpAt:   timestamppb.New(time.Now().Add(time.Hour * 12)),
				Creator: creator,
			},
			err: true,
		},
		{
			name: "user dont have role creator",
			in: &promos.CreatePromo{
				Name:    uuid.NewString(),
				Uses:    -1,
				Amount:  -10,
				ExpAt:   timestamppb.New(time.Now().Add(time.Hour * 12)),
				Creator: user,
			},
			err: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Creating promo
			out, err := client.Create(context.TODO(), tt.in)
			if (err != nil) != tt.err {
				t.Fail()
			}

			if out != nil {
				promoIds = append(promoIds, out.PromoCode.Id)
			}
		})
	}

}

func Use(t *testing.T) {
	var userId int64
	var promoId int64

	t.Cleanup(
		func() {
			// Delete history activation of promo by user
			if err := repo.DeleteActivatePromoFromHistory(context.TODO(), pool, &promos.PromoId{Id: promoId}); err != nil {
				t.Fatal(err)
			}

			// Delete testing data from table Promos
			if err := clearPromos([]int64{promoId}); err != nil {
				t.Fatal(err)
			}

			// Delete testing data from table Users
			if err := clearUsers([]int64{userId}); err != nil {
				t.Fatal(err)
			}
		},
	)

	// Creating creator user
	userId, err := serviceUsers.Create(context.TODO(), 3)
	if err != nil {
		t.Fatal(err)
	}

	promoFailure, err := client.Create(context.TODO(), &promos.CreatePromo{
		Name:    uuid.NewString(),
		Uses:    1,
		ExpAt:   timestamppb.New(time.Now().Add(time.Hour * 12)),
		Creator: userId,
	})
	if err != nil {
		t.Fatal(err)
	}
	promoId = promoFailure.PromoCode.Id

	var tests = []struct {
		name string
		in   *promos.PromoUserId
		err  bool
	}{
		{
			name: "common",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: false,
		},
		{
			name: "already activated",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: true,
		},
		{
			name: "not in stock",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: true,
		},
		{
			name: "expired",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := client.Use(context.TODO(), tt.in); (err != nil) != tt.err {
				t.Fail()
			}
		})

		// Add uses, only for test
		if tt.name == "common" {
			if _, err := client.AddUses(context.TODO(), &promos.AddUsesIn{PromoId: promoId, Uses: 1}); err != nil {
				t.Fatal(err)
			}
		}

		// Delete history activation promo by user, only for tests
		if tt.name == "already activated" {
			if err := repo.DeleteActivatePromoFromHistory(context.TODO(), pool, &promos.PromoId{Id: promoId}); err != nil {
				t.Fatal(err)
			}
			if _, err := client.AddUses(context.TODO(), &promos.AddUsesIn{PromoId: promoId, Uses: -1}); err != nil {
				t.Fatal(err)
			}
		}

		// Expire promo, only for tests
		if tt.name == "not in stock" {
			if _, err := client.AddTime(context.TODO(), &promos.AddTimeIn{PromoId: promoId, ExpAt: timestamppb.New(time.Date(2000, time.April, 16, 10, 0, 0, 0, time.UTC))}); err != nil {
				t.Fatal(err)
			}
			if _, err := client.AddUses(context.TODO(), &promos.AddUsesIn{PromoId: promoId, Uses: 1}); err != nil {
				t.Fatal(err)
			}
		}
	}

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

func clearPromos(promoIds []int64) error {
	for _, promoId := range promoIds {
		if _, err := client.Delete(context.TODO(), &promos.PromoId{Id: promoId}); err != nil {
			return err
		}
	}

	return nil
}
