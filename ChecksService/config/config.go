package config

import (
	Err "checks/pkg/errors"
	"checks/pkg/grpc/conn"
	"checks/pkg/grpc/server"
	"checks/pkg/hasher"
	"checks/pkg/postgres"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server server.ServerConfig
	DB     postgres.DBConfig
	Conn   conn.Config
	Hash   hasher.Config
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
	srvUsersHost, srvUsersPort := os.Getenv("SERVICE_USERS_HOST"), os.Getenv("SERVICE_USERS_PORT")
	if srvUsersHost == "" || srvUsersPort == "" {
		return Config{}, fmt.Errorf("%s %s", Err.ErrMissingEnviroment, "for connection to service users")
	}

	// Config warns
	salt := os.Getenv("HASH_SALT")
	if salt == "" {
		return Config{}, fmt.Errorf("%s %s", Err.ErrMissingEnviroment, "for hasher")
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
				Host: srvUsersHost,
				Port: srvUsersPort,
			},
		},
		Hash: hasher.Config{
			Salt: salt,
		},
	}, nil
}
