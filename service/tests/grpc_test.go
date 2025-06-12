package test_grpc

import (
	"context"
	"errors"
	"fmt"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"
	"promos/tests/mock"
	"testing"
	"time"

	Err "promos/pkg/errors"
	"promos/pkg/postgres"

	"github.com/google/uuid"
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
var mockpool *mock.MockPool
var logger *logrus.Logger

// DONT EDIT PROMO
var promo = &promos.CreatePromo{
	Name:     uuid.New().String(),
	Amount:   0,
	Currency: 0,
	Uses:     2,
	ExpAt:    timestamppb.New(time.Now().Add(time.Second * 1000)),
}

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
	mockpool, err = mock.MockPoolUp(mock.Config{User: cfg.DB.User, Password: cfg.DB.Password, Name: cfg.DB.Database})
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
		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := mockpool.MockPoolDown(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ALL TESTING", func(t *testing.T) {
		var userId1, userId2, userId3, userId4, promoId int64

		t.Cleanup(func() {
			userIds := []int64{userId1, userId2, userId3, userId4}
			promoIds := []int64{promoId}

			// Delete testing data from table UserToPromo
			if err := clearUserToPromo(userIds, promoId); err != nil {
				t.Fatal(err)
			}
			// Delete testing data from table Promos
			if err := clearPromos(promoIds); err != nil {
				t.Fatal(err)
			}
			// Delete testing data from table Users
			if err := clearUsers(userIds); err != nil {
				t.Fatal(err)
			}
		})

		// Adding user1
		userId1, err := serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}
		promo.Creator = userId1

		// Adding promo
		promoFailure, err := client.Create(context.TODO(), promo)
		if err != nil {
			t.Fail()
		}
		promoId = promoFailure.PromoCode.Id

		// TEST GetById
		if _, err := client.GetById(context.TODO(), &promos.PromoId{Id: promoId}); err != nil {
			logger.Warn("error test GetById promo fail", err)
			t.Fail()
		}

		// TEST GetByName
		if _, err := client.GetByName(context.TODO(), &promos.PromoName{Name: promoFailure.PromoCode.Name}); err != nil {
			logger.Warn("error test GetByName promo fail", err)
			t.Fail()
		}

		// TEST Use
		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId1}); err != nil {
			logger.Warn("error test Use promo fail", err)
			t.Fail()
		}

		// TEST Use. Second attemp use promo by one user. Want error, promo is already activated by user
		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId1}); errors.Is(err, Err.ErrPromoAlreadyActivated) {
			logger.Warn("error test second Use promo by one user fail", err)
			t.Fail()
		}

		// Adding user2
		userId2, err = serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}
		promo.Creator = userId2

		// TEST Use. Second activation of promo.
		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId2}); err != nil {
			logger.Warn("error test second Use promo fail", err)
			t.Fail()
		}

		// Adding user3
		userId3, err = serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}
		promo.Creator = userId3

		// TEST Use. Third activation of promo. Want error, promo not in stock
		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId2}); errors.Is(err, Err.ErrPromoNotInStock) {
			logger.Warn("error test third Use promo, promo not in stock, fail", err)
			t.Fail()
		}

		// Sleeping 1 second, for testing promo expiring
		time.Sleep(time.Second * 1)

		// Adding user4
		userId4, err = serviceUsers.Create(context.TODO())
		if err != nil {
			t.Fail()
		}
		promo.Creator = userId4

		// TEST Use. Fourth activation of promo. Want error, promo is expired
		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId4}); errors.Is(err, Err.ErrPromoExpired) {
			logger.Warn("error test fourth Use promo, promo expired, fail", err)
			t.Fail()
		}

	})
}

func clearUserToPromo(userIds []int64, promoId int64) error {
	for _, userId := range userIds {
		if err := serviceUsers.ClearHistory(context.TODO(), userId, promoId); err != nil {
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

func clearPromos(promoIds []int64) error {
	for _, promoId := range promoIds {
		if _, err := client.DeleteById(context.TODO(), &promos.PromoId{Id: promoId}); err != nil {
			return err
		}
	}

	return nil
}
