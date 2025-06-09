package test_grpc

import (
	"context"
	"fmt"
	"promos/config"
	"promos/internal/models/promos"
	"promos/internal/repository"
	"promos/internal/transport/grpc/conn"
	service "promos/internal/transport/grpc/handlers"
	"promos/internal/transport/grpc/server"
	"testing"

	"promos/pkg/postgres"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var srv *server.Server
var client promos.PromosClient

func init() {

	// Logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

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
		Conn: conn.Config{
			CfgSrvUsers: conn.ConfigServiceUsers{
				Host: "localhost",
				Port: "50124",
			},
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
	grpcSrv := grpc.NewServer()

	// Connect to other services
	clientServices, err := conn.NewClientsServices(cfg.Conn)
	if err != nil {
		panic(err)
	}

	// Creating repository
	repo := repository.NewRepository(clientServices)

	// Register promo service
	service := service.NewServicePromos(repo, pool, logger)
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

	t.Run("Test Create | Use | GetById | GetByName | Delete", func(t *testing.T) {
		createOut, err := client.Create(context.TODO(), promo)
		if err != nil {
			t.Fatal("create: ", err)
		}

		getByIdOut, err := client.GetById(context.TODO(), &promos.PromoId{Id: createOut.PromoCode.Id})
		if err != nil {
			t.Fatal("getById: ", err)
		}

		getByNameOut, err := client.GetByName(context.TODO(), &promos.PromoName{Name: createOut.PromoCode.Name})
		if err != nil {
			t.Fatal("getByName: ", err)
		}

		if getByIdOut.PromoCode.Id != getByNameOut.PromoCode.Id {
			t.Fatal(getByIdOut, getByNameOut)
		}

		if _, err := client.Use(context.TODO(), &promos.PromoUserId{PromoId: getByIdOut.PromoCode.Id, UserId: 4}); err != nil {
			t.Fatal("use: ", err)
		}

		if _, err := client.DeleteById(context.TODO(), &promos.PromoId{Id: createOut.PromoCode.Id}); err != nil {
			t.Fatal("delete: ", err)
		}
	})

}
