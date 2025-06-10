package repository

import (
	"context"
	"promos/internal/models/users"

	"google.golang.org/grpc"
)

type UserService interface {
	SendTransaction(ctx context.Context, in *users.TransactionRequest, opts ...grpc.CallOption) (*users.TransactionResponse, error)
}

type Repository struct {
	UserService
}

func NewRepository(userService UserService) *Repository {
	return &Repository{UserService: userService}
}
