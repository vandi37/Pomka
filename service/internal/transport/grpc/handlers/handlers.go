package service

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	repeatible "promos/pkg/utils"

	"github.com/jackc/pgx/v5"
)

func (sp *ServicePromos) Create(ctx context.Context, in *promos.CreatePromoIn) (out *promos.CreatePromoOut, err error) {
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {
		out, err = sp.repo.Create(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return out, nil
}

func (sp *ServicePromos) Delete(ctx context.Context, in *promos.PromoName) (out *common.Response, err error) {
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {
		out, err = sp.repo.Delete(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return out, nil
}

func (sp *ServicePromos) Use(ctx context.Context, in *promos.PromoName) (*users.TransactionResponse, error) {
	return nil, nil
}
