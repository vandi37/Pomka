package service

import (
	"context"
	"errors"
	"fmt"
	"warns/pkg/models/common"
	"warns/pkg/models/users"
	"warns/pkg/models/warns"
	repeatible "warns/pkg/utils"

	Err "warns/pkg/errors"

	"github.com/jackc/pgx/v5"
)

func (sw *ServiceWarns) Warn(ctx context.Context, in *warns.WarnCreate) (warnsFailure *warns.WarnFailure, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {

		// Get info about moderator
		user, err := sw.users.GetUser(ctx, &users.Id{Id: in.ModerId})
		if err != nil {
			errors.Join(Err.ErrServiceUsers, err)
		}

		// Check moderator role
		if b, err := sw.repo.IsUserModerator(ctx, user); !b || err != nil {
			return errors.Join(err)
		}

		// Create warn for this user
		warnsFailure.Warn, err = sw.repo.CreateWarn(ctx, tx, in)
		if err != nil {
			return errors.Join(Err.ErrCreateWarn, err)
		}

		// Send transaction to service users
		if _, err := sw.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Warn,
		}); err != nil {
			return errors.Join(Err.ErrSendTransaction, err)
		}

		// Check count of warns for this user
		cntWarns, err := sw.repo.GetCountOfActiveWarns(ctx, tx, &users.Id{Id: in.UserId})
		if err != nil {
			return errors.Join(Err.ErrCountActiveWarns, err)
		}

		// Pass, if user dont have too many warns
		if sw.cfg.WarnsBeforeBan != cntWarns {
			return nil
		}

		// Insert ban into Bans
		banReason := fmt.Sprintf("user got %d warns", cntWarns)
		if _, err := sw.repo.CreateBan(ctx, tx, &warns.BanCreate{
			UserId:  in.UserId,
			ModerId: in.ModerId,
			Reason:  &banReason,
		}); err != nil {
			return errors.Join(Err.ErrCreateBan, err)
		}

		// Send transaction to service users
		if _, err := sw.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Ban,
		}); err != nil {
			return errors.Join(Err.ErrSendTransaction, err)
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return warnsFailure, nil
}

func (sw *ServiceWarns) AllUnWarn(ctx context.Context, in *users.Id) (*common.Response, error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		if err := sw.repo.MakeWarnsInActive(ctx, tx, in); err != nil {
			return errors.Join(Err.ErrMakeWarnsInActive, err)
		}

		// Send transaction to service users
		if _, err := sw.users.SendTransaction(ctx, &users.TransactionRequest{
			Type: common.TransactionType_InActiveWarn,
		}); err != nil {
			return errors.Join(Err.ErrSendTransaction, err)
		}

		return nil

	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

func (sw *ServiceWarns) LastUnWarn(ctx context.Context, in *users.Id) (*common.Response, error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

func (sw *ServiceWarns) Ban(ctx context.Context, in *warns.BanCreate) (*warns.BanFailure, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {

		// Insert ban into bans
		if _, err := sw.repo.CreateBan(ctx, tx, in); err != nil {
			return errors.Join(Err.ErrCreateBan, err)
		}

		// Send transaction to service users
		if _, err := sw.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Block,
		}); err != nil {
			return errors.Join(Err.ErrSendTransaction, err)
		}

		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

func (sw *ServiceWarns) Unban(ctx context.Context, in *users.Id) (*common.Response, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {

		// Remove ban
		if err := sw.repo.MakeBanInActive(ctx, tx, in); err != nil {
			return errors.Join(Err.ErrMakeBansInActive, err)
		}

		// Send transaction to service users
		if _, err := sw.users.SendTransaction(ctx, &users.TransactionRequest{
			Type: common.TransactionType_User,
		}); err != nil {
			return errors.Join(Err.ErrSendTransaction, err)
		}

		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

func (sw *ServiceWarns) GetHistoryWarns(ctx context.Context, in *users.Id) (warnsFailure *warns.AllWarnsFailure, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		warnsFailure.Warns, err = sw.repo.GetWarns(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return warnsFailure, nil
}

func (sw *ServiceWarns) GetHistoryBans(ctx context.Context, in *users.Id) (bansFailure *warns.AllBansFailure, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		bansFailure.Bans, err = sw.repo.GetBans(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return bansFailure, nil
}
func (sw *ServiceWarns) GetActiveWarns(ctx context.Context, in *users.Id) (warnsFailure *warns.AllWarnsFailure, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		warnsFailure.Warns, err = sw.repo.GetWarns(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return warnsFailure, nil
}

func (sw *ServiceWarns) GetActiveBan(ctx context.Context, in *users.Id) (banFailure *warns.BanFailure, err error) {

	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return banFailure, nil
}

func (sw *ServiceWarns) GetCountOfActiveWarns(ctx context.Context, in *users.Id) (*warns.CountOfActiveWarns, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}
