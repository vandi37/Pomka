package service

import (
	"context"
	"promos/internal/models/promos"
	"promos/internal/models/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServicePromos struct {
	repo RepositoryPromos
	db   *pgxpool.Pool
	promos.UnimplementedPromosServer
}

type RepositoryPromos interface {
	CreatePromo(ctx context.Context, tx pgx.Tx, in *promos.CreatePromo) (out *promos.PromoCode, err error)
	DeletePromoById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (err error)
	DeletePromoByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (err error)
	GetPromoById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (out *promos.PromoCode, err error)
	GetPromoByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (out *promos.PromoCode, err error)

	PromoIsValid(in *promos.PromoCode) (b bool, err error)
	PromoIsAlreadyActivated(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) (b bool, err error)

	ActivatePromo(ctx context.Context, in *promos.PromoCode, userId int64) (out *users.TransactionResponse, err error)
	AddActivatePromoToHistory(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) (err error)
	DeleteActivatePromoFromHistory(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) (err error)
	DecrementPromoUses(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (err error)
}

func NewServicePromos(repo RepositoryPromos, db *pgxpool.Pool) *ServicePromos {
	return &ServicePromos{repo: repo, db: db}
}
