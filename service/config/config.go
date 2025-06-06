package config

import (
	"fmt"
	"os"
	"promos/internal/transport/grpc/server"
	"promos/pkg/postgres"
	"strconv"
)

type Config struct {
	Server server.ServerConfig
	DB     postgres.DBConfig
}

func NewConfig() (Config, error) {

	// Config server
	srvNet, srvPort :=
		os.Getenv("SERVER_NETWORK"),
		os.Getenv("SERVER_PORT")
	if srvNet == "" {
		srvNet = "tcp"
	}
	if srvPort == "" {
		srvPort = "50123"
	}

	// Config db
	dbHost, dbPort, dbUser, dbPassword, dbName, dbMaxAtmps, dbDelayAtmps :=
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_MAX_ATMPS"),
		os.Getenv("DB_DELAY_ATMPS_S")
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "mAz0H1zm"
	}
	if dbName == "" {
		dbName = "postgres"
	}
	if dbMaxAtmps == "" {
		dbMaxAtmps = "5"
	}
	if dbDelayAtmps == "" {
		dbDelayAtmps = "5"
	}

	dbMaxAtmpsInt, err1 := strconv.Atoi(dbMaxAtmps)
	dbDelayAtmpsInt, err2 := strconv.Atoi(dbDelayAtmps)
	if err1 != nil || err2 != nil {
		return Config{}, fmt.Errorf("config: NewConfig: error wrong enviroment param")
	}

	return Config{
		Server: server.ServerConfig{
			Network: srvNet,
			Port:    srvPort,
		},
		DB: postgres.DBConfig{
			Host:        dbHost,
			Port:        dbPort,
			User:        dbUser,
			Password:    dbPassword,
			Database:    dbName,
			MaxAtmps:    dbMaxAtmpsInt,
			DelayAtmpsS: dbDelayAtmpsInt,
		},
	}, nil
}
