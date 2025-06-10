package repository

import (
	"context"
	"log"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	Err "promos/pkg/errors"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Create(ctx context.Context, tx pgx.Tx, in *promos.CreatePromo) (*promos.PromoFailure, error) {
	q := `INSERT INTO Promos (Id, Name, Currency, Amount, Uses, Creator, ExpAt)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`

	expAt := in.ExpAt.AsTime().Format("2006-01-02 15:04:05")
	id := int64(uuid.New().ID())
	_, err := tx.Exec(ctx, q,
		id,
		in.Name,
		in.Currency,
		in.Amount,
		in.Uses,
		in.Creator,
		expAt)
	if err != nil {
		return nil, Err.ErrExecQuery
	}

	out := &promos.PromoFailure{
		PromoCode: &promos.PromoCode{
			Id:        id,
			Name:      in.Name,
			Currency:  in.Currency,
			Amount:    in.Amount,
			Uses:      in.Uses,
			Creator:   in.Creator,
			ExpAt:     in.ExpAt,
			CreatedAt: timestamppb.Now()},
		Failure: nil}
	return out, nil
}

func (r *Repository) DeleteById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (*common.Response, error) {
	q := `DELETE FROM UserToPromo WHERE PromoId = $1`

	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
		log.Println(err)
		return nil, Err.ErrExecQuery
	}

	q = `DELETE FROM Promos WHERE Id = $1`

	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
		log.Println(err)
		return nil, Err.ErrExecQuery
	}

	return nil, nil
}

func (r *Repository) DeleteByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*common.Response, error) {
	q := `DELETE FROM Promos WHERE Name = $1`

	if _, err := tx.Exec(ctx, q, in.Name); err != nil {
		return nil, Err.ErrExecQuery
	}

	return nil, nil
}

func (r *Repository) GetPromoById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (*promos.PromoCode, error) {
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Id = $1`
	row := tx.QueryRow(ctx, q, in.Id)

	if err := row.Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator, nil, nil); err != nil {
		return nil, Err.ErrIncorrectData
	}
	return out, nil
}

func (r *Repository) GetPromoByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*promos.PromoCode, error) {
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Name = $1`
	row := tx.QueryRow(ctx, q, in.Name)

	if err := row.Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator, nil, nil); err != nil {
		return nil, Err.ErrIncorrectData
	}

	return out, nil
}

func (r *Repository) Activate(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoCode,
	userId int64) (*users.TransactionResponse, error) {

	out, err := r.UserService.SendTransaction(ctx, &users.TransactionRequest{
		Sender:   nil,
		Receiver: &users.UserTransaction{UserId: userId, Amount: in.Amount, Currency: in.Currency},
	})
	if err != nil {
		return nil, Err.ErrServiceUsers
	}

	return out, nil
}

func (r *Repository) IsValid(in *promos.PromoCode) error {
	if float64(in.ExpAt.AsTime().Unix()) > float64(time.Now().Unix()) {
		return Err.ErrPromoExpired
	}

	if in.Uses == 0 {
		return Err.ErrPromoNotInStock
	}

	return nil
}

func (r *Repository) DecrementUses(ctx context.Context, tx pgx.Tx, in *promos.PromoId) error {
	if in.Id == -1 {
		return nil
	}

	q := `UPDATE Promos SET Uses = Uses-1 WHERE Id = $1`

	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
		return Err.ErrExecQuery
	}

	return nil
}

func (r *Repository) AddUserToPromo(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) error {
	q := `INSERT INTO UserToPromo (UserId, PromoId, ActivatedAt) VALUES ($1, $2, $3)`

	actAt := time.Now().Format("2006-01-02 15:04:05")
	if _, err := tx.Exec(ctx, q, in.UserId, in.PromoId, actAt); err != nil {
		return Err.ErrExecQuery
	}

	return nil
}

func (r *Repository) IsAlreadyActivated(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) error {
	var activated bool

	q := `SELECT EXISTS(SELECT * FROM UserToPromo WHERE UserId = $1 AND PromoId = $2)`

	row := tx.QueryRow(ctx, q, in.UserId, in.PromoId)
	if err := row.Scan(&activated); err != nil {
		log.Println(err)
		return Err.ErrIncorrectData
	}

	return nil
}
