package repository

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	Err "promos/pkg/errors"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(ctx context.Context, tx pgx.Tx, in *promos.CreatePromoIn) (*promos.CreatePromoOut, error) {
	q := `INSERT INTO promos (Name, Value, Creator, Currency, ExpAt)
    VALUES ($1, $2, $3, $4, $5)`

	expAt := in.ExpAt.AsTime().Format("2006-01-02 15:04:05")

	_, err := tx.Exec(ctx, q,
		in.Name,
		in.Value,
		in.Creator,
		in.Currency,
		expAt)

	if err != nil {
		return nil, Err.ErrExecQuery
	}

	out := &promos.CreatePromoOut{
		PromoCode: &promos.PromoCode{
			Name:      in.Name,
			Value:     in.Value,
			Creator:   in.Creator,
			Currency:  in.Currency,
			ExpAt:     in.ExpAt,
			CreatedAt: timestamppb.Now()}}

	return out, nil
}

func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*common.Response, error) {
	q := `DELETE FROM promos WHERE Name=$1`

	if _, err := tx.Exec(ctx, q, in.Name); err != nil {
		return nil, Err.ErrExecQuery
	}

	return nil, nil
}

func (r *Repository) Use(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*users.TransactionResponse, error) {
	return nil, nil
}
