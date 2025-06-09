package service

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	repeatible "promos/pkg/utils"

	"github.com/jackc/pgx/v5"
)

func (sp *ServicePromos) Create(ctx context.Context, in *promos.CreatePromo) (out *promos.PromoFailure, err error) {
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

func (sp *ServicePromos) DeleteById(ctx context.Context, in *promos.PromoId) (out *common.Response, err error) {
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {
		out, err = sp.repo.DeleteById(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return out, nil
}

func (sp *ServicePromos) DeleteByName(ctx context.Context, in *promos.PromoName) (out *common.Response, err error) {
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {
		out, err = sp.repo.DeleteByName(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return out, nil
}

func (sp *ServicePromos) Use(ctx context.Context, in *promos.PromoUserId) (out *users.TransactionResponse, err error) {
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Get PromoCode
		promo, err := sp.repo.GetPromoById(ctx, tx, &promos.PromoId{Id: in.PromoId})
		if err != nil {
			return err
		}

		// Check valid promo
		if err := sp.repo.IsValid(promo); err != nil {
			return err
		}

		// Send query to user service
		out, err = sp.repo.Activate(ctx, tx, promo, in.UserId)
		if err != nil {
			return err
		}

		// Increment uses promo
		if err := sp.repo.DecrementUses(ctx, tx, &promos.PromoId{Id: in.PromoId}); err != nil {
			return err
		}

		// Add activation promo to history
		if err := sp.repo.AddUserToPromo(ctx, tx, in); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return out, nil
}

func (sp *ServicePromos) GetById(ctx context.Context, in *promos.PromoId) (out *promos.PromoFailure, err error) {
	var promo *promos.PromoCode

	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {
		promo, err = sp.repo.GetPromoById(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return &promos.PromoFailure{PromoCode: promo, Failure: nil}, nil
}

func (sp *ServicePromos) GetByName(ctx context.Context, in *promos.PromoName) (out *promos.PromoFailure, err error) {
	var promo *promos.PromoCode

	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {
		promo, err = sp.repo.GetPromoByName(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return &promos.PromoFailure{PromoCode: promo, Failure: nil}, nil
}
