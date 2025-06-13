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
	"time"

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
		logger.Debug("Cleanup TestMain start")

		// Stoping gRPC server
		srv.Stop()

		// Stoping docker postgres
		if err := mockpool.MockPoolDown(); err != nil {
			t.Fatal(err)
		}

		logger.Debug("Cleanup TestMain done")
	})

	t.Run("CREATE AND DELETE PROMO", CreateDelete)
	t.Run("COMMON USE PROMO", CommonUse)
	t.Run("USE PROMO ALREADY ACTIVATED", UsePromoAlreadyActivated)
	t.Run("USE PROMO NOT IN STOCK", UsePromoAlreadyActivated)
}

func CreateDelete(t *testing.T) {
	var promoId, userId int64
	t.Cleanup(func() {
		logger.Debug("Cleanup TestCreateDelete start")

		userIds := []int64{userId}
		promoIds := []int64{promoId}

		// Delete testing data from table Promos
		if err := clearPromos(promoIds); err != nil {
			t.Fatal(err)
		}
		// Delete testing data from table Users
		if err := clearUsers(userIds); err != nil {
			t.Fatal(err)
		}

		logger.Debug("Cleanup TestCreateDelete done")
	})

	var createIn = &promos.CreatePromo{
		Name:  uuid.New().String(),
		ExpAt: timestamppb.New(time.Now().Add(time.Second * 1)),
	}

	// Creating user
	userId, err := serviceUsers.Create(context.TODO())
	if err != nil {
		logger.Warn("Create user fail")
		t.Fail()
	}
	createIn.Creator = userId

	// Creating promo
	createOut, err := client.Create(context.TODO(), createIn)
	if err != nil {
		logger.Warn("Create promo fail")
		t.Fail()
	}
	promoId = createOut.PromoCode.Id
}

func CommonUse(t *testing.T) {
	var promoId, userId int64
	t.Cleanup(func() {
		logger.Debug("Cleanup CommonUse start")

		userIds := []int64{userId}
		promoIds := []int64{promoId}

		// Delete testing data from table Promos
		if err := clearPromos(promoIds); err != nil {
			t.Fatal(err)
		}
		// Delete testing data from table Users
		if err := clearUsers(userIds); err != nil {
			t.Fatal(err)
		}

		logger.Debug("Cleanup CommonUse done")
	})

	var createIn = &promos.CreatePromo{
		Name:  uuid.New().String(),
		Uses:  1,
		ExpAt: timestamppb.New(time.Now().Add(time.Second * 1)),
	}

	// Creating user
	userId, err := serviceUsers.Create(context.TODO())
	if err != nil {
		logger.Warn("Create user fail")
		t.Fail()
	}
	createIn.Creator = userId

	// Creating promo
	createOut, err := client.Create(context.TODO(), createIn)
	if err != nil {
		logger.Warn("Create promo fail")
		t.Fail()
	}
	promoId = createOut.PromoCode.Id

	if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId}); err != nil {
		logger.Warn("Use promo fail")
	}
}

func UsePromoAlreadyActivated(t *testing.T) {
	var promoId, userId int64
	t.Cleanup(func() {
		logger.Debug("Cleanup TestUsePromoAlreadyActivated start")

		userIds := []int64{userId}
		promoIds := []int64{promoId}

		// Delete testing data from table Promos
		if err := clearPromos(promoIds); err != nil {
			t.Fatal(err)
		}
		// Delete testing data from table Users
		if err := clearUsers(userIds); err != nil {
			t.Fatal(err)
		}

		logger.Debug("Cleanup TestUsePromoAlreadyActivated done")
	})

	var createIn = &promos.CreatePromo{
		Name:  uuid.New().String(),
		Uses:  1,
		ExpAt: timestamppb.New(time.Now().Add(time.Second * 1)),
	}

	// Creating user
	userId, err := serviceUsers.Create(context.TODO())
	if err != nil {
		logger.Warn("Create user fail")
		t.Fail()
	}
	createIn.Creator = userId

	// Creating promo
	createOut, err := client.Create(context.TODO(), createIn)
	if err != nil {
		logger.Warn("Create promo fail")
		t.Fail()
	}
	promoId = createOut.PromoCode.Id

	if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId}); err != nil {
		logger.Warn("Use promo fail")
	}
	if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userId}); err == nil {
		logger.Warn("Use promo fail")
	}
}

func UsePromoNotInStock(t *testing.T) {
	var promoId, userIdFirst, userIdSecond int64
	t.Cleanup(func() {
		logger.Debug("Cleanup TestPromoNotInStock start")

		userIds := []int64{userIdFirst, userIdSecond}
		promoIds := []int64{promoId}

		// Delete testing data from table Promos
		if err := clearPromos(promoIds); err != nil {
			t.Fatal(err)
		}
		// Delete testing data from table Users
		if err := clearUsers(userIds); err != nil {
			t.Fatal(err)
		}

		logger.Debug("Cleanup TestPromoNotInStock done")
	})

	var createIn = &promos.CreatePromo{
		Name:  uuid.New().String(),
		Uses:  1,
		ExpAt: timestamppb.New(time.Now().Add(time.Second * 1)),
	}

	// Creating user
	userIdFirst, err := serviceUsers.Create(context.TODO())
	if err != nil {
		logger.Warn("Create user fail")
		t.Fail()
	}
	createIn.Creator = userIdFirst

	// Creating promo
	createOut, err := client.Create(context.TODO(), createIn)
	if err != nil {
		logger.Warn("Create promo fail")
		t.Fail()
	}
	promoId = createOut.PromoCode.Id

	if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userIdFirst}); err != nil {
		logger.Warn("Use promo fail")
	}

	// Creating user
	userIdSecond, err = serviceUsers.Create(context.TODO())
	if err != nil {
		logger.Warn("Create user fail")
		t.Fail()
	}
	createIn.Creator = userIdSecond

	if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: promoId, UserId: userIdSecond}); err != nil {
		logger.Warn("Use promo fail")
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
