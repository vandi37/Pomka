package service

import (
	"context"
	"postgres"
	"protobuf/checks"
	"protobuf/users"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type Config struct {
	WarnsBeforeBan int
}

type UserService interface {
	SendTransaction(ctx context.Context, in *users.TransactionRequest, opts ...grpc.CallOption) (*users.TransactionResponse, error)
	GetUser(ctx context.Context, in *users.Id, opts ...grpc.CallOption) (*users.User, error)
}

type ServiceChecks struct {
	RepositoryChecks
	db *pgxpool.Pool
	UserService
	checks.UnimplementedChecksServer
}

type RepositoryChecks interface {
	CreateCheck(ctx context.Context, db postgres.DB, in *checks.CheckCreate) (out *checks.Check, err error)
	RemoveCheck(ctx context.Context, db postgres.DB, in *checks.CheckId) (err error)
	GetUsersCheck(ctx context.Context, db postgres.DB, in *users.Id) (out *checks.AllChecks, err error)
	GetCheckByKey(ctx context.Context, db postgres.DB, key string) (out *checks.Check, err error)
}

func NewServiceChecks(repo RepositoryChecks, db *pgxpool.Pool, users UserService) *ServiceChecks {
	return &ServiceChecks{RepositoryChecks: repo, db: db, UserService: users}
}
