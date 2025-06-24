package promos_test

import (
	"context"
	"fmt"
	"promos/config"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	Err "promos/pkg/errors"
	migrations "promos/pkg/goose"
	"promos/pkg/grpc/server"
	"promos/pkg/models/promos"
	"promos/tests/mock"
	"testing"
	"time"

	"promos/pkg/postgres"

	"promos/pkg/logger"

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
var dockerpostgres *mock.DockerPool
var repo *repository.Repository
var pool *pgxpool.Pool
var log *logrus.Logger

func TestMain(m *testing.M) {
	defer func() {

		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := dockerpostgres.PostgresDown(); err != nil {
			log.Fatal(err)
		}

	}()

	// Setup logger
	log = logger.NewLogger()

	// Load enviroment
	if err := godotenv.Load("config.env"); err != nil {
		log.WithField("ERROR", err).Panic("SETUP APP")
	}
	log.WithField("MSG", "Succecs loading enviroment file").Debug("SETUP APP")

	// Cofiguration
	cfg, err := config.NewConfig()
	if err != nil {
		log.WithField("ERROR", err).Panic("SETUP APP")
	}
	log.WithField("MSG", "Succecs loading configuration for app").Debug("SETUP APP")

	// Up docker postgres
	dockerpostgres, err = mock.PostgresUp(mock.Config{User: cfg.DB.User, Password: cfg.DB.Password, Name: cfg.DB.Database})
	if err != nil {
		log.WithField("ERROR", err).Panic("SETUP APP")
	}
	log.WithField("MSG", "Succecs up postgres docker container").Debug("SETUP APP")

	// Connecting to postgres
	pool, err = postgres.NewPool(context.TODO(), cfg.DB)
	if err != nil {
		log.WithField("ERROR", err).Panic("SETUP APP")
	}
	log.WithField("MSG", "Succecs connect to postgres").Debug("SETUP APP")

	// Run migrations
	if err := migrations.Up(context.TODO(), pool); err != nil {
		log.WithField("ERROR", err).Panic("SETUP APP")
	}
	log.WithField("MSG", "Succecs run migrations").Debug("SETUP APP")

	// gRPC server
	grpcSrv := grpc.NewServer(grpc.UnaryInterceptor(server.NewServerLogger(log).LoggingUnaryInterceptor))

	// Creating mock service users
	serviceUsers = mock.NewMockServiceUsers(pool)

	// Creating repository
	repo = repository.NewRepository()

	// Register promo service
	service := service.NewServicePromos(repo, pool, serviceUsers)
	promos.RegisterPromosServer(grpcSrv, service)

	// Run server
	srv = server.NewServer(grpcSrv)
	go func() {
		if err := srv.Run(cfg.Server); err != nil {
			log.WithField("ERROR", err).Panic("SETUP APP")
		}
	}()
	log.WithField("MSG", fmt.Sprintf("Running server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("SETUP APP")

	// Connection to server
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%s", cfg.Server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.WithField("ERROR", err).Panic("SETUP APP")
	}
	client = promos.NewPromosClient(conn)
	log.WithField("MSG", fmt.Sprintf("Succecs connection to server on %s:%s", cfg.Server.Network, cfg.Server.Port)).Debug("SETUP APP")

	m.Run()
}

func TestCreateDelete(t *testing.T) {
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

	// Creating user with role creator
	creator, err := serviceUsers.Create(context.TODO(), 3)
	if err != nil {
		t.Fatal(err)
	}

	// Creating user with role user
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

func TestUse(t *testing.T) {
	var userId int64
	var promoId int64

	t.Cleanup(
		func() {

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
		err  error
	}{
		{
			name: "common",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: nil,
		},
		{
			name: "already activated",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: Err.ErrPromoAlreadyActivated,
		},
		{
			name: "not in stock",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: Err.ErrPromoNotInStock,
		},
		{
			name: "expired",
			in: &promos.PromoUserId{
				PromoId: promoId,
				UserId:  userId,
			},
			err: Err.ErrPromoExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := client.Use(context.TODO(), tt.in); !(fmt.Errorf("rpc error: code = Unknown desc = %s", tt.err) != err) {
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
