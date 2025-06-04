package config

import (
	"fmt"
	"os"
	"promos/internal/transport/grpc/server"
	"promos/pkg/postgres"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server server.ServerConfig
	DB     postgres.DBConfig
}

func NewConfig() (Config, error) {
	godotenv.Load("./env")

	// Config server
	srvNet, srvPort :=
		os.Getenv("SERVER_NETWORK"),
		os.Getenv("SERVER_PORT")
	if srvNet == "" || srvPort == "" {
		return Config{}, fmt.Errorf("config: NewConfig: error missing server env params")
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
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbMaxAtmps == "" || dbDelayAtmps == "" {
		return Config{}, fmt.Errorf("config: NewConfig: error missing db env params")
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
