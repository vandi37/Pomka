package repository

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, in *promos.CreatePromoIn) (*promos.CreatePromoOut, error) {
	return nil, nil
}

func (r *Repository) Delete(ctx context.Context, in *promos.PromoName) (*common.Response, error) {
	return nil, nil
}

func (r *Repository) Use(ctx context.Context, in *promos.PromoName) (*users.TransactionResponse, error) {
	return nil, nil
}
