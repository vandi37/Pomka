package service

import (
	"context"
	"protobuf/common"
	"protobuf/market"
	"utils"

	"github.com/jackc/pgx/v5"
)

func (s *ServiceMarket) Sell(ctx context.Context, in *market.OfferCreate) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return nil, nil
}

func (s *ServiceMarket) Buy(ctx context.Context, in *market.TypeName) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return nil, nil
}

func (s *ServiceMarket) GetOffers(ctx context.Context, in *market.OffersSorting) (*market.OffersFailure, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	var offersFailure = new(market.OffersFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return &market.OffersFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return offersFailure, nil
}

func (s *ServiceMarket) RemoveFromSell(ctx context.Context, in *market.TypeName) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return &common.Response{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return nil, nil
}
