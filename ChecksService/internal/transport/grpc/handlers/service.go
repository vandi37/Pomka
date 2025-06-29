package service

import (
	"context"
	"postgres"
	"protobuf/checks"
	"protobuf/users"

	"conn"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	WarnsBeforeBan int
}

type ServiceChecks struct {
	RepositoryChecks
	db *pgxpool.Pool
	conn.UserService
	checks.UnimplementedChecksServer
}

type RepositoryChecks interface {
	CreateCheck(ctx context.Context, db postgres.DB, in *checks.CheckCreate) (out *checks.Check, err error)
	RemoveCheck(ctx context.Context, db postgres.DB, in *checks.CheckId) (err error)
	GetUsersCheck(ctx context.Context, db postgres.DB, in *users.Id) (out *checks.AllChecks, err error)
	GetCheckByKey(ctx context.Context, db postgres.DB, key string) (out *checks.Check, err error)
}

func NewServiceChecks(repo RepositoryChecks, db *pgxpool.Pool, users conn.UserService) *ServiceChecks {
	return &ServiceChecks{RepositoryChecks: repo, db: db, UserService: users}
}
