package service

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
)

func (sp *ServicePromos) Create(ctx context.Context, in *promos.CreatePromoIn) (*promos.CreatePromoOut, error) {
	return nil, nil
}

func (sp *ServicePromos) Delete(ctx context.Context, in *promos.PromoName) (*common.Response, error) {
	return nil, nil
}

func (sp *ServicePromos) Use(ctx context.Context, in *promos.PromoName) (*users.TransactionResponse, error) {
	return nil, nil
}
