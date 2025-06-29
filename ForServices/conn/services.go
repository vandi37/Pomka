package conn

import (
	"context"
	"protobuf/users"

	"google.golang.org/grpc"
)

type UserService interface {
	SendTransaction(ctx context.Context, in *users.TransactionRequest, opts ...grpc.CallOption) (*users.TransactionResponse, error)
	GetUser(ctx context.Context, in *users.Id, opts ...grpc.CallOption) (*users.User, error)
}

// temp mock
type CmdsService interface {
	ChangeOwner()
}
