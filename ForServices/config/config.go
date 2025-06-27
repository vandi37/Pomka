package config

import (
	"os"
	"strconv"

	"conn"
	e "errorspomka"
	"postgres"
	"server"
)

type Config struct {
	Server  server.ServerConfig
	DB      postgres.Config
	Conn    conn.Config
	Storage Storage
}

func NewConfig() (Config, error) {

	// Config server
	srvNet, srvPort :=
		os.Getenv("SERVER_NETWORK"),
		os.Getenv("SERVER_PORT")

	if srvNet == "" || srvPort == "" {
		return Config{}, e.ErrMissingEnviroment
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
		return Config{}, e.ErrMissingEnviroment
	}
	dbMaxAtmpsInt, err1 := strconv.Atoi(dbMaxAtmps)
	dbDelayAtmpsInt, err2 := strconv.Atoi(dbDelayAtmps)
	if err1 != nil || err2 != nil {
		return Config{}, e.ErrMissingEnviroment
	}

	// Config connection to service users
	srvUsersHost, srvUsersPort :=
		os.Getenv("SERVICE_USERS_HOST"),
		os.Getenv("SERVICE_USERS_PORT")
	if srvUsersHost == "" || srvUsersPort == "" {
		return Config{}, e.ErrMissingEnviroment
	}

	// Config hasher
	salt := os.Getenv("HASH_SALT")
	if salt == "" {
		return Config{}, e.ErrMissingEnviroment
	}

	// Config warns
	warnsBeforeBan := os.Getenv("WARNS_BEFORE_BAN")
	warnsBeforeBanInt, err := strconv.Atoi(warnsBeforeBan)
	if err != nil || warnsBeforeBanInt <= 0 {
		return Config{}, e.ErrMissingEnviroment
	}

	return Config{
		Server: server.ServerConfig{
			Network: srvNet,
			Port:    srvPort,
		},
		DB: postgres.Config{
			Host:        dbHost,
			Port:        dbPort,
			User:        dbUser,
			Password:    dbPassword,
			Database:    dbName,
			MaxAtmps:    dbMaxAtmpsInt,
			DelayAtmpsS: dbDelayAtmpsInt,
		},
		Conn: conn.Config{
			ConfigServiceUsers: conn.ConfigServiceUsers{
				Host: srvUsersHost,
				Port: srvUsersPort,
			},
		},
		Storage: Storage{
			HashSalt:       salt,
			WarnsBeforeBan: warnsBeforeBanInt,
		},
	}, nil
}
