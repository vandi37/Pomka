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
	Create(ctx context.Context, tx pgx.Tx, in *promos.CreatePromoIn) (*promos.CreatePromoOut, error)
	Delete(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*common.Response, error)
	Use(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*users.TransactionResponse, error)
}

func NewServicePromos(repo RepositoryPromos, db *pgxpool.Pool, log *logrus.Logger) *ServicePromos {
	return &ServicePromos{repo: repo, db: db, log: log}
}
