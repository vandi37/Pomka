package repository

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
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
