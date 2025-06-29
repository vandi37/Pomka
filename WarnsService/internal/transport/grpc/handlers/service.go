package service

import (
	"context"
	"protobuf/users"
	"protobuf/warns"

	"postgres"

	"conn"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	WarnsBeforeBan int
}

type ServiceWarns struct {
	repo  RepositoryWarns
	db    *pgxpool.Pool
	cfg   Config
	users conn.UserService
	warns.UnimplementedWarnsServer
}

type RepositoryWarns interface {
	CreateWarn(ctx context.Context, db postgres.DB, in *warns.ModerUserReason) (warn *warns.Warn, err error)
	CreateBan(ctx context.Context, db postgres.DB, in *warns.ModerUserReason) (warn *warns.Ban, err error)
	GetWarns(ctx context.Context, db postgres.DB, in *users.Id) (warn *warns.AllWarns, err error)
	GetBans(ctx context.Context, db postgres.DB, in *users.Id) (warn *warns.AllBans, err error)
	GetActiveWarns(ctx context.Context, db postgres.DB, in *users.Id) (warn *warns.AllWarns, err error)
	GetActiveBan(ctx context.Context, db postgres.DB, in *users.Id) (warn *warns.Ban, err error)
	IsUserModerator(ctx context.Context, in *users.User) (b bool, err error)
	MakeWarnsInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error)
	MakeLastWarnInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error)
	MakeBanInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error)
	GetCountOfActiveWarns(ctx context.Context, db postgres.DB, in *users.Id) (*warns.CountOfActiveWarns, error)
	IsAlreadyBanned(ctx context.Context, db postgres.DB, in *users.Id) (b bool, err error)
}

func NewServiceWarns(repo RepositoryWarns, db *pgxpool.Pool, cfg Config, users conn.UserService) *ServiceWarns {
	return &ServiceWarns{repo: repo, db: db, cfg: cfg, users: users}
}
