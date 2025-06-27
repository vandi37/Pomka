package service

import (
	"context"
	"protobuf/checks"
	"protobuf/common"
	"protobuf/users"

	"utils"

	"github.com/jackc/pgx/v5"
)

func (s *ServiceChecks) Create(ctx context.Context, in *checks.CheckCreate) (checkFailure *checks.CheckFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	checkFailure = new(checks.CheckFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Create check
		checkFailure.Check, err = s.CreateCheck(ctx, tx, in)
		if err != nil {
			return err
		}

		// Send transation to service users
		if _, err := s.UserService.SendTransaction(
			ctx, &users.TransactionRequest{
				Sender: &users.UserTransaction{
					UserId:   in.Creator,
					Amount:   in.Amount,
					Currency: in.Currency,
				},
				Type: common.TransactionType_CreateCheck,
			},
		); err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &checks.CheckFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return checkFailure, nil
}

func (s *ServiceChecks) Remove(ctx context.Context, in *checks.CheckId) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Remove check
		if err := s.RemoveCheck(ctx, tx, in); err != nil {
			return err
		}

		// Send transaction to service users
		if _, err := s.UserService.SendTransaction(
			ctx, &users.TransactionRequest{Type: common.TransactionType_DeleteCheck}); err != nil {
			return err
		}

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

func (s *ServiceChecks) Use(ctx context.Context, in *checks.CheckUse) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get check
		check, err := s.GetCheckByKey(ctx, tx, in.Key)
		if err != nil {
			codeError = common.ErrorCode_CheckNotValid
			return err
		}

		// Send transaction to service users
		if _, err := s.UserService.SendTransaction(
			ctx, &users.TransactionRequest{
				Sender: &users.UserTransaction{
					UserId:   in.UserId,
					Amount:   check.Amount,
					Currency: check.Currency,
				},
				Type: common.TransactionType_CreateCheck,
			},
		); err != nil {
			return err
		}

		if err := s.RemoveCheck(ctx, tx, &checks.CheckId{Id: check.Id}); err != nil {
			return err
		}

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

func (s *ServiceChecks) GetUserChecks(ctx context.Context, in *users.Id) (allChecksFailure *checks.AllChecksFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	allChecksFailure = new(checks.AllChecksFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get checks user
		allChecksFailure.AllChecks, err = s.GetUsersCheck(ctx, tx, in)
		if err != nil {
			return err
		}
		return nil

	}); errTx != nil {
		return &checks.AllChecksFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return allChecksFailure, nil
}
