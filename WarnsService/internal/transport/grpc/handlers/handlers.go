package service

import (
	"context"
	"errors"
	"fmt"
	"protobuf/common"
	"protobuf/users"
	"protobuf/warns"
	"utils"

	e "errorspomka"

	"github.com/jackc/pgx/v5"
)

func (s *ServiceWarns) Warn(ctx context.Context, in *warns.ModerUserReason) (warnsFailure *warns.WarnFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	warnsFailure = new(warns.WarnFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Check user already banned
		if b, err := s.repo.IsAlreadyBanned(ctx, tx, &users.Id{Id: in.UserId}); b || err == e.ErrExecQuery {
			codeError = common.ErrorCode_UserAlreadyBanned
			return err
		}

		// Get info about moderator
		user, err := s.users.GetUser(ctx, &users.Id{Id: in.ModerId})
		if err != nil {
			errors.Join(e.ErrServiceUsers, err)
		}

		// Check moderator role
		if b, err := s.repo.IsUserModerator(ctx, user); !b || err != nil {
			codeError = common.ErrorCode_UserBadRole
			return errors.Join(err)
		}

		// Create warn for this user
		warnsFailure.Warn, err = s.repo.CreateWarn(ctx, tx, in)
		if err != nil {
			return errors.Join(e.ErrCreateWarn, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Warn,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
		}

		// Check count of warns for this user
		cntWarns, err := s.repo.GetCountOfActiveWarns(ctx, tx, &users.Id{Id: in.UserId})
		if err != nil {
			return errors.Join(e.ErrCountActiveWarns, err)
		}

		// Pass, if user dont have too many warns
		if s.cfg.WarnsBeforeBan != int(cntWarns.CountWarns) {
			return nil
		}

		// Make warns for this user inactive
		if err := s.repo.MakeWarnsInActive(ctx, tx, &users.Id{Id: in.UserId}); err != nil {
			return errors.Join(e.ErrMakeWarnsInActive, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_InActiveWarn,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
		}

		// Insert ban into Bans
		banReason := fmt.Sprintf("user got %d warns", cntWarns.CountWarns)
		if _, err := s.repo.CreateBan(ctx, tx, &warns.ModerUserReason{
			UserId:  in.UserId,
			ModerId: in.ModerId,
			Reason:  &banReason,
		}); err != nil {
			return errors.Join(e.ErrCreateBan, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Ban,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
		}

		return nil

	}); errTx != nil {
		return &warns.WarnFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return warnsFailure, nil
}

func (s *ServiceWarns) AllUnWarn(ctx context.Context, in *warns.ModerUserReason) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get info about moderator
		user, err := s.users.GetUser(ctx, &users.Id{Id: in.ModerId})
		if err != nil {
			errors.Join(e.ErrServiceUsers, err)
		}

		// Check moderator role
		if b, err := s.repo.IsUserModerator(ctx, user); !b || err != nil {
			codeError = common.ErrorCode_UserBadRole
			return errors.Join(err)
		}

		if err := s.repo.MakeWarnsInActive(ctx, tx, &users.Id{Id: in.UserId}); err != nil {
			return errors.Join(e.ErrMakeWarnsInActive, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_InActiveWarn,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
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

func (s *ServiceWarns) LastUnWarn(ctx context.Context, in *warns.ModerUserReason) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get info about moderator
		user, err := s.users.GetUser(ctx, &users.Id{Id: in.ModerId})
		if err != nil {
			errors.Join(e.ErrServiceUsers, err)
		}

		// Check moderator role
		if b, err := s.repo.IsUserModerator(ctx, user); !b || err != nil {
			codeError = common.ErrorCode_UserBadRole
			return errors.Join(err)
		}

		if err := s.repo.MakeLastWarnInActive(ctx, tx, &users.Id{Id: in.UserId}); err != nil {
			return errors.Join(e.ErrMakeWarnsInActive, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_InActiveWarn,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
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

func (s *ServiceWarns) Ban(ctx context.Context, in *warns.ModerUserReason) (banFailure *warns.BanFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	banFailure = new(warns.BanFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get info about moderator
		user, err := s.users.GetUser(ctx, &users.Id{Id: in.ModerId})
		if err != nil {
			errors.Join(e.ErrServiceUsers, err)
		}

		// Check moderator role
		if b, err := s.repo.IsUserModerator(ctx, user); !b || err != nil {
			codeError = common.ErrorCode_UserBadRole
			return errors.Join(err)
		}

		// Check user already banned
		if b, err := s.repo.IsAlreadyBanned(ctx, tx, &users.Id{Id: in.UserId}); b || err != nil {
			return err
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Block,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
		}

		// Insert ban into bans
		banFailure.Ban, err = s.repo.CreateBan(ctx, tx, in)
		if err != nil {
			return errors.Join(e.ErrCreateBan, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_Ban,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
		}

		return nil

	}); errTx != nil {
		return &warns.BanFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return banFailure, nil
}

func (s *ServiceWarns) Unban(ctx context.Context, in *warns.ModerUserReason) (*common.Response, error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {

		// Get info about moderator
		user, err := s.users.GetUser(ctx, &users.Id{Id: in.ModerId})
		if err != nil {
			errors.Join(e.ErrServiceUsers, err)
		}

		// Check moderator role
		if b, err := s.repo.IsUserModerator(ctx, user); !b || err != nil {
			codeError = common.ErrorCode_UserBadRole
			return errors.Join(err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_User,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
		}

		// Remove ban
		if err := s.repo.MakeBanInActive(ctx, tx, &users.Id{Id: in.UserId}); err != nil {
			return errors.Join(e.ErrMakeBansInActive, err)
		}

		// Send transaction to service users
		if _, err := s.users.SendTransaction(ctx, &users.TransactionRequest{
			Sender:   &users.UserTransaction{UserId: in.ModerId},
			Receiver: &users.UserTransaction{UserId: in.UserId},
			Type:     common.TransactionType_InActiveBan,
		}); err != nil {
			return errors.Join(e.ErrSendTransaction, err)
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

func (s *ServiceWarns) GetHistoryWarns(ctx context.Context, in *users.Id) (warnsFailure *warns.AllWarnsFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	warnsFailure = new(warns.AllWarnsFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		warnsFailure.Warns, err = s.repo.GetWarns(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return &warns.AllWarnsFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return warnsFailure, nil
}

func (s *ServiceWarns) GetHistoryBans(ctx context.Context, in *users.Id) (bansFailure *warns.AllBansFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	bansFailure = new(warns.AllBansFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		bansFailure.Bans, err = s.repo.GetBans(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return &warns.AllBansFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return bansFailure, nil
}

func (s *ServiceWarns) GetActiveWarns(ctx context.Context, in *users.Id) (warnsFailure *warns.AllWarnsFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	warnsFailure = new(warns.AllWarnsFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		warnsFailure.Warns, err = s.repo.GetActiveWarns(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return &warns.AllWarnsFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return warnsFailure, nil
}

func (s *ServiceWarns) GetActiveBan(ctx context.Context, in *users.Id) (banFailure *warns.BanFailure, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden
	banFailure = new(warns.BanFailure)

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		banFailure.Ban, err = s.repo.GetActiveBan(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return &warns.BanFailure{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return banFailure, nil
}

func (s *ServiceWarns) GetCountOfActiveWarns(ctx context.Context, in *users.Id) (count *warns.CountOfActiveWarns, err error) {
	var codeError common.ErrorCode = common.ErrorCode_Forbidden

	// Run in transaction
	if errTx := utils.RunInTx(s.db, ctx, func(tx pgx.Tx) error {
		count, err = s.repo.GetCountOfActiveWarns(ctx, tx, in)
		if err != nil {
			return err
		}

		return nil

	}); errTx != nil {
		return &warns.CountOfActiveWarns{
			Failure: &common.Failure{
				Code: codeError,
				Details: map[string]string{
					"ERROR": errTx.Error(),
				},
			},
		}, errTx
	}

	return count, nil
}
