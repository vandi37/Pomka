package service

import (
	"context"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	repo "promos/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServicePromos struct {
	repo RepositoryPromos
	db   *pgxpool.Pool
	promos.UnimplementedPromosServer
}

type RepositoryPromos interface {
	CreatePromo(ctx context.Context, db repo.DB, in *promos.CreatePromo) (out *promos.PromoCode, err error)
	DeletePromoById(ctx context.Context, db repo.DB, in *promos.PromoId) (err error)
	DeletePromoByName(ctx context.Context, db repo.DB, in *promos.PromoName) (err error)
	GetPromoById(ctx context.Context, db repo.DB, in *promos.PromoId) (out *promos.PromoCode, err error)
	GetPromoByName(ctx context.Context, db repo.DB, in *promos.PromoName) (out *promos.PromoCode, err error)

	PromoIsExpired(in *promos.PromoCode) (b bool, err error)
	PromoIsNotInStock(in *promos.PromoCode) (b bool, err error)
	PromoIsAlreadyActivated(ctx context.Context, db repo.DB, in *promos.PromoUserId) (b bool, err error)
	CreatorIsOwner(ctx context.Context, userId int64) (b bool, err error)

	ActivatePromo(ctx context.Context, in *promos.PromoCode, userId int64) (out *users.TransactionResponse, err error)
	AddActivatePromoToHistory(ctx context.Context, db repo.DB, in *promos.PromoUserId) (err error)
	DeleteActivatePromoFromHistory(ctx context.Context, db repo.DB, in *promos.PromoId) (err error)
	DecrementPromoUses(ctx context.Context, db repo.DB, in *promos.PromoId) (err error)

	AddTime(ctx context.Context, db repo.DB, in *promos.AddTimeIn) (err error)
	AddUses(ctx context.Context, db repo.DB, in *promos.AddUsesIn) (err error)
}

func NewServicePromos(repo RepositoryPromos, db *pgxpool.Pool) *ServicePromos {
	return &ServicePromos{repo: repo, db: db}
}
