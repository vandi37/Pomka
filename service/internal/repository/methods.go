package repository

import (
	"context"
	"errors"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	Err "promos/pkg/errors"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Insert promo to table promos
func (r *Repository) CreatePromo(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.CreatePromo) (*promos.PromoCode, error) {

	var expiredAt, createdAt interface{}
	var out = new(promos.PromoCode)

	q := `INSERT INTO Promos (Name, Currency, Amount, Uses, Creator, ExpAt)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING Id, Name, Currency, Amount, Uses, Creator, ExpAt, CreatedAt`
	expAt := in.ExpAt.AsTime().Format("2006-01-02 15:04:05")

	err := tx.QueryRow(ctx, q,
		in.Name,
		in.Currency,
		in.Amount,
		in.Uses,
		in.Creator,
		expAt).Scan(&out.Id, &out.Name, &out.Currency, &out.Amount, &out.Uses, &out.Creator, &expiredAt, &createdAt)

	if err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return nil, Err.ErrExecQuery
	}

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt.(time.Time)), timestamppb.New(createdAt.(time.Time))
	return out, nil
}

// Delete promo from table promos by ID
func (r *Repository) DeletePromoById(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoId) error {
	q := `DELETE FROM Promos WHERE Id = $1`

	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return Err.ErrExecQuery
	}

	return nil
}

// Delete promo from table promos by Name
func (r *Repository) DeletePromoByName(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoName) error {
	q := `DELETE FROM Promos WHERE Name = $1`

	if _, err := tx.Exec(ctx, q, in.Name); err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return Err.ErrExecQuery
	}

	return nil
}

// Get promo from table promos by ID
func (r *Repository) GetPromoById(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoId) (*promos.PromoCode, error) {
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Id = $1`
	row := tx.QueryRow(ctx, q, in.Id)

	err := row.Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator, nil, nil)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			r.logger.Warn(errors.Join(Err.ErrMissingPromoId, err))
			return nil, Err.ErrMissingPromoId
		default:
			r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
			return nil, Err.ErrExecQuery
		}
	}

	return out, nil
}

// Get promo from table promos by Name
func (r *Repository) GetPromoByName(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoName) (*promos.PromoCode, error) {
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Name = $1`
	row := tx.QueryRow(ctx, q, in.Name)

	err := row.Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator, nil, nil)

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			r.logger.Warn(errors.Join(Err.ErrMissingPromoName, err))
			return nil, Err.ErrMissingPromoName
		default:
			r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
			return nil, Err.ErrExecQuery
		}
	}

	return out, nil
}

// Send transaction to service users
func (r *Repository) ActivatePromo(
	ctx context.Context,
	in *promos.PromoCode,
	userId int64) (*users.TransactionResponse, error) {

	out, err := r.UserService.SendTransaction(ctx, &users.TransactionRequest{
		Sender:   nil,
		Receiver: &users.UserTransaction{UserId: userId, Amount: in.Amount, Currency: in.Currency},
	})
	if err != nil {
		r.logger.Warn(errors.Join(Err.ErrServiceUsers, err))
		return nil, Err.ErrServiceUsers
	}

	return out, nil
}

// Update table promos, decrement uses of promo.
func (r *Repository) DecrementPromoUses(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoId) (err error) {
	if in.Id == -1 {
		return nil
	}

	q := `UPDATE Promos SET Uses = Uses-1 WHERE Id = $1`

	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return Err.ErrExecQuery
	}

	return nil
}

// Insert activation of promo to table UserToPromo
func (r *Repository) AddActivatePromoToHistory(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoUserId) (err error) {
	q := `INSERT INTO UserToPromo (UserId, PromoId, ActivatedAt) VALUES ($1, $2, $3)`

	actAt := time.Now().Format("2006-01-02 15:04:05")
	if _, err := tx.Exec(ctx, q, in.UserId, in.PromoId, actAt); err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return Err.ErrExecQuery
	}

	return nil
}

// Delete activation of promo from table UserToPromo
func (r *Repository) DeleteActivatePromoFromHistory(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoUserId) (err error) {

	q := `DELETE FROM UserToPromo WHERE UserId=$1 AND PromoId=$2`
	if _, err := tx.Exec(ctx, q, in.UserId, in.PromoId); err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return Err.ErrExecQuery
	}

	return nil
}

// If promo been activated by user, return true. If promo not activated by user, return false.
func (r *Repository) PromoIsAlreadyActivated(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoUserId) (b bool, err error) {

	var activated = new(bool)

	q := `SELECT EXISTS(SELECT * FROM UserToPromo WHERE UserId = $1 AND PromoId = $2)`

	row := tx.QueryRow(ctx, q, in.UserId, in.PromoId)
	if err := row.Scan(&activated); err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return false, Err.ErrExecQuery
	}

	return *activated, nil
}

// If promo valid, return true. If promo is expired/uses=0, return false.
func (r *Repository) PromoIsValid(in *promos.PromoCode) (b bool, err error) {
	if float64(in.ExpAt.AsTime().Unix()) > float64(time.Now().Unix()) {
		return false, Err.ErrPromoExpired
	}

	if in.Uses == 0 {
		return false, Err.ErrPromoNotInStock
	}

	return true, nil
}
