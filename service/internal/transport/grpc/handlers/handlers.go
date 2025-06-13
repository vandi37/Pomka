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
	var promo *promos.PromoCode

	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Creating promo
		promo, err = sp.repo.CreatePromo(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &promos.PromoFailure{PromoCode: nil, Failure: nil}, errTx
	}

	return &promos.PromoFailure{PromoCode: promo, Failure: nil}, nil
}

func (sp *ServicePromos) Delete(ctx context.Context, in *promos.PromoId) (out *common.Response, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Deleting history
		if err := sp.repo.DeleteActivatePromoFromHistory(ctx, tx, in); err != nil {
			return err
		}

		// Deleting promo
		if err := sp.repo.DeletePromoById(ctx, tx, in); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &common.Response{Failure: nil}, errTx
	}

	return &common.Response{Failure: nil}, nil
}

func (sp *ServicePromos) Use(ctx context.Context, in *promos.PromoUserId) (out *users.TransactionResponse, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Query to db for get promo
		promo, err := sp.repo.GetPromoById(ctx, tx, &promos.PromoId{Id: in.PromoId})
		if err != nil {
			return err
		}

		// Check creator is owner or not
		if b, err := sp.repo.CreatorIsOwner(ctx, promo); err != nil || !b {
			return err
		}

		// Check promo is expired or not
		if b, err := sp.repo.PromoIsExpired(promo); err != nil || !b {
			return err
		}

		// Check promo is in stock or not
		if b, err := sp.repo.PromoIsNotInStock(promo); err != nil || !b {
			return err
		}

		// Query to db for check activation promo from user
		if b, err := sp.repo.PromoIsAlreadyActivated(ctx, tx, in); err != nil || b {
			return err
		}

		// Query to serviceUsers for give currency
		out, err = sp.repo.ActivatePromo(ctx, promo, in.UserId)
		if err != nil {
			return err
		}

		// Query to db for decremenet uses of promo
		if err := sp.repo.DecrementPromoUses(ctx, tx, &promos.PromoId{Id: in.PromoId}); err != nil {
			return err
		}

		// Query to db for adding promo activation in history
		if err := sp.repo.AddActivatePromoToHistory(ctx, tx, in); err != nil {
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

	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Get promo
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

	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Get promo
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

func (sp *ServicePromos) AddTime(ctx context.Context, in *promos.AddTimeIn) (*common.Response, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Add time for promo
		if err := sp.repo.AddTime(ctx, tx, in); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

func (sp *ServicePromos) AddUses(ctx context.Context, in *promos.AddUsesIn) (*common.Response, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sp.db, ctx, func(tx pgx.Tx) error {

		// Add uses for promo
		if err := sp.repo.AddUses(ctx, tx, in); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}
