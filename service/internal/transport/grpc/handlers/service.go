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
	Create(ctx context.Context, tx pgx.Tx, in *promos.CreatePromo) (out *promos.PromoFailure, err error)
	DeleteById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (out *common.Response, err error)
	DeleteByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (out *common.Response, err error)
	Activate(ctx context.Context, tx pgx.Tx, in *promos.PromoCode, userId int64) (out *users.TransactionResponse, err error)
	GetPromoById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (out *promos.PromoCode, err error)
	GetPromoByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (out *promos.PromoCode, err error)
	IsValid(in *promos.PromoCode) (err error)
	DecrementUses(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (err error)
	AddUserToPromo(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) (err error)
	IsAlreadyActivated(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) (err error)
}

func NewServicePromos(repo RepositoryPromos, db *pgxpool.Pool, log *logrus.Logger) *ServicePromos {
	return &ServicePromos{repo: repo, db: db, log: log}
}
