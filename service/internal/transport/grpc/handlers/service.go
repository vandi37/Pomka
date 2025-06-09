package service

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type ServicePromos struct {
	log  *logrus.Logger
	repo RepositoryPromos
	db   *pgxpool.Pool
	promos.UnimplementedPromosServer
}

type RepositoryPromos interface {
	Create(ctx context.Context, tx pgx.Tx, in *promos.CreatePromo) (*promos.PromoFailure, error)
	DeleteById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (*common.Response, error)
	DeleteByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*common.Response, error)
	Activate(ctx context.Context, tx pgx.Tx, in *promos.PromoCode, userId int64) (*users.TransactionResponse, error)
	GetPromoById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (*promos.PromoCode, error)
	GetPromoByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*promos.PromoCode, error)
	IsValid(in *promos.PromoCode) error
	DecrementUses(ctx context.Context, tx pgx.Tx, in *promos.PromoId) error
	AddUserToPromo(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) error
}

func NewServicePromos(repo RepositoryPromos, db *pgxpool.Pool, log *logrus.Logger) *ServicePromos {
	return &ServicePromos{repo: repo, db: db, log: log}
}
