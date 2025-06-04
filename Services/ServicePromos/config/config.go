package config

import (
	"fmt"
	"os"
	"promos/internal/transport/grpc/server"

	"github.com/joho/godotenv"
)

type Config struct {
	Server server.ServerConfig
}

func parseEnv(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return err
	}

	return nil
}

func NewConfig() (Config, error) {
	if err := parseEnv("./config/.env"); err != nil {
		return Config{}, fmt.Errorf("config: NewConfig: %s", err)
	}

	srvNet, srvPort := os.Getenv("ServerNetwork"), os.Getenv("ServerPort")
	if srvNet == "" || srvPort == "" {
		return Config{}, fmt.Errorf("config: NewConfig: error missing env params")
	}

	return Config{
		Server: server.ServerConfig{
			Network: srvNet,
			Port:    srvPort,
		},
	}, nil
}
