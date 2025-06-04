package service

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
)

type ServicePromos struct {
	repo RepositoryPromos
	promos.UnimplementedPromosServer
}

type RepositoryPromos interface {
	Create(ctx context.Context, in *promos.CreatePromoIn) (*promos.CreatePromoOut, error)
	Delete(ctx context.Context, in *promos.PromoName) (*common.Response, error)
	Use(ctx context.Context, in *promos.PromoName) (*users.TransactionResponse, error)
}

func NewServicePromos(repo RepositoryPromos) *ServicePromos {
	return &ServicePromos{repo: repo}
}
