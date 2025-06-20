package config

import (
	"fmt"
	"os"
	"strconv"
	service "warns/internal/transport/grpc/handlers"
	Err "warns/pkg/errors"
	"warns/pkg/grpc/conn"
	"warns/pkg/grpc/server"
	"warns/pkg/postgres"
)

type Config struct {
	Server server.ServerConfig
	DB     postgres.DBConfig
	Conn   conn.Config
	Warns  service.Config
}

func NewConfig() (Config, error) {

	// Config server
	srvNet, srvPort :=
		os.Getenv("SERVER_NETWORK"),
		os.Getenv("SERVER_PORT")

	if srvNet == "" || srvPort == "" {
		return Config{}, fmt.Errorf("%s %s", Err.ErrMissingEnviroment, "for server")
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
		return Config{}, fmt.Errorf("%s %s", Err.ErrMissingEnviroment, "for db")
	}
	dbMaxAtmpsInt, err1 := strconv.Atoi(dbMaxAtmps)
	dbDelayAtmpsInt, err2 := strconv.Atoi(dbDelayAtmps)
	if err1 != nil || err2 != nil {
		return Config{}, Err.ErrMissingEnviroment
	}

	// Config connection to service users
	SrvUsersHost, SrvUsersPort := os.Getenv("SERVICE_USERS_HOST"), os.Getenv("SERVICE_USERS_PORT")
	if SrvUsersHost == "" || SrvUsersPort == "" {
		return Config{}, fmt.Errorf("%s %s", Err.ErrMissingEnviroment, "for connection to service users")
	}

	// Config warns
	warnsBeforeBan := os.Getenv("WARNS_BEFORE_BAN")
	warnsBeforeBanInt, err := strconv.Atoi(warnsBeforeBan)
	if err != nil || warnsBeforeBanInt <= 0 {
		return Config{}, nil
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
		Conn: conn.Config{
			CfgSrvUsers: conn.ConfigServiceUsers{
				Host: SrvUsersHost,
				Port: SrvUsersPort,
			},
		},
		Warns: service.Config{
			WarnsBeforeBan: warnsBeforeBanInt,
		},
	}, nil
}
