package conn

import (
	"fmt"

	"protobuf/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientsServices struct {
	conn *grpc.ClientConn
	users.UsersClient
}

func NewClientsServices(cfg Config) (*ClientsServices, error) {

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", cfg.ConfigServiceUsers.Host, cfg.ConfigServiceUsers.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientUsers := users.NewUsersClient(conn)
	return &ClientsServices{
		conn,
		clientUsers,
	}, nil
}

func (c ClientsServices) Close() error {
	return c.conn.Close()
}
