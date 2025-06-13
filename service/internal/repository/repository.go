package repository

import (
	"context"
	"promos/internal/models/users"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UserService interface {
	SendTransaction(ctx context.Context, in *users.TransactionRequest, opts ...grpc.CallOption) (*users.TransactionResponse, error)
	GetUser(ctx context.Context, in *users.Id, opts ...grpc.CallOption) (*users.User, error)
}

type Repository struct {
	logger *logrus.Logger
	UserService
}

func NewRepository(userService UserService, logger *logrus.Logger) *Repository {
	return &Repository{UserService: userService, logger: logger}
}
