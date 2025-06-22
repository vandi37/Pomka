package service

import (
	"context"
	"promos/pkg/models/common"
	"promos/pkg/models/promos"
	"promos/pkg/models/users"
	repeatible "promos/pkg/utils"

	"github.com/jackc/pgx/v5"
)

func (s *ServicePromos) Create(ctx context.Context, in *promos.CreatePromo) (promoFailure *promos.PromoFailure, err error) {
	promoFailure = new(promos.PromoFailure)

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Query to service users for get information about creator
		user, err := s.users.GetUser(ctx, &users.Id{Id: in.Creator})
		if err != nil {
			return err
		}

		// Check creator is owner or not
		if b, err := s.repo.CreatorIsOwner(ctx, user); err != nil || !b {
			return err
		}

		// Creating promo
		promoFailure.PromoCode, err = s.repo.CreatePromo(ctx, tx, in)
		if err != nil {
			return err
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender: &users.UserTransaction{UserId: in.Creator},
			Type:   common.TransactionType_CreatePromoCode,
		}); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &promos.PromoFailure{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			}}, errTx
	}

	return promoFailure, nil
}

func (s *ServicePromos) Delete(ctx context.Context, in *promos.PromoId) (*common.Response, error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Deleting history
		if err := s.repo.DeleteActivatePromoFromHistory(ctx, tx, in); err != nil {
			return err
		}

		// Deleting promo
		if err := s.repo.DeletePromoById(ctx, tx, in); err != nil {
			return err
		}

		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Type: common.TransactionType_DeletePromoCode,
		}); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			}}, errTx
	}

	return nil, nil
}

func (s *ServicePromos) Use(ctx context.Context, in *promos.PromoUserId) (*common.Response, error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Query to db for get promo
		promo, err := s.repo.GetPromoById(ctx, tx, &promos.PromoId{Id: in.PromoId})
		if err != nil {
			return err
		}

		// Check promo is expired or not
		if b, err := s.repo.PromoIsExpired(promo); err != nil || !b {
			return err
		}

		// Check promo is in stock or not
		if b, err := s.repo.PromoIsNotInStock(promo); err != nil || !b {
			return err
		}

		// Query to db for check activation promo from user
		if b, err := s.repo.PromoIsAlreadyActivated(ctx, tx, in); err != nil || b {
			return err
		}

		// Query to db for decrement uses of promo
		if err := s.repo.DecrementPromoUses(ctx, tx, &promos.PromoId{Id: in.PromoId}); err != nil {
			return err
		}

		// Send transaction to service users
		if _, err = s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender: &users.UserTransaction{UserId: in.UserId},
			Type:   common.TransactionType_DecrementUsesPromo,
		}); err != nil {
			return err
		}

		// Query to db for adding promo activation in history
		if err := s.repo.AddActivatePromoToHistory(ctx, tx, in); err != nil {
			return err
		}

		// Send transaction to service users
		if _, err = s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender: &users.UserTransaction{UserId: in.UserId},
			Type:   common.TransactionType_AddActivationPromoCodeToHistory,
		}); err != nil {
			return err
		}

		// Send transaction to service users
		if _, err = s.users.SendTransaction(ctx, &users.TransactionRequest{
			Receiver: &users.UserTransaction{UserId: in.UserId, Amount: promo.Amount, Currency: promo.Currency},
			Type:     common.TransactionType_ActivatePromoCode,
		}); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			}}, errTx
	}

	return nil, nil
}

func (s *ServicePromos) GetById(ctx context.Context, in *promos.PromoId) (promoFailure *promos.PromoFailure, err error) {
	promoFailure = new(promos.PromoFailure)

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get promo
		promoFailure.PromoCode, err = s.repo.GetPromoById(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &promos.PromoFailure{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return promoFailure, nil
}

func (s *ServicePromos) GetByName(ctx context.Context, in *promos.PromoName) (promoFailure *promos.PromoFailure, err error) {
	promoFailure = new(promos.PromoFailure)

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get promo
		promoFailure.PromoCode, err = s.repo.GetPromoByName(ctx, tx, in)

		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &promos.PromoFailure{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			}}, errTx
	}

	return promoFailure, nil
}

func (s *ServicePromos) AddTime(ctx context.Context, in *promos.AddTimeIn) (*common.Response, error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Add time for promo
		if err := s.repo.AddTime(ctx, tx, in); err != nil {
			return err
		}

		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Type: common.TransactionType_AddTimeForPromo,
		}); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			}}, errTx
	}

	return nil, nil
}

func (s *ServicePromos) AddUses(ctx context.Context, in *promos.AddUsesIn) (*common.Response, error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Add uses for promo
		if err := s.repo.AddUses(ctx, tx, in); err != nil {
			return err
		}

		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Type: common.TransactionType_AddUsesForPromo,
		}); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: common.ErrorCode_Promos,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			}}, errTx
	}

	return nil, nil
}
