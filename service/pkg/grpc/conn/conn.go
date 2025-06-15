package conn

import (
	"fmt"
	"warns/pkg/models/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientsServices struct {
	users.UsersClient
}

func NewClientsServices(cfg Config) (*ClientsServices, error) {

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.CfgSrvUsers.Host, cfg.CfgSrvUsers.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientUsers := users.NewUsersClient(conn)
	return &ClientsServices{
		clientUsers,
	}, nil
}
