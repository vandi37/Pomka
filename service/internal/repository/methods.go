package repository

import (
	"context"
	"errors"
	"promos/internal/models/common"
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

	if in.Uses < 0 && in.Uses != -1 {
		return nil, Err.ErrValueUses
	}
	if float64(time.Now().Unix()) > float64(in.ExpAt.AsTime().Unix()) {
		return nil, Err.ErrExpAt
	}

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
		expAt).Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator,
		&expiredAt,
		&createdAt)

	if err != nil {
		r.logger.Warn(errors.Join(Err.ErrExecQuery, err))
		return nil, Err.ErrExecQuery
	}

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt.(time.Time)), timestamppb.New(createdAt.(time.Time))

	r.logger.Debugf("creating promo: %d", out.Id)
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

	r.logger.Debugf("deleting promo: %d", in.Id)
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

	r.logger.Debugf("deleting promo: %s", in.Name)
	return nil
}

// Get promo from table promos by ID
func (r *Repository) GetPromoById(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoId) (*promos.PromoCode, error) {

	var expiredAt, createdAt interface{}
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Id = $1`
	row := tx.QueryRow(ctx, q, in.Id)

	err := row.Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator,
		&expiredAt,
		&createdAt)

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

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt.(time.Time)), timestamppb.New(createdAt.(time.Time))
	return out, nil
}

// Get promo from table promos by Name
func (r *Repository) GetPromoByName(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoName) (*promos.PromoCode, error) {

	var expiredAt, createdAt interface{}
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Name = $1`
	row := tx.QueryRow(ctx, q, in.Name)

	err := row.Scan(
		&out.Id,
		&out.Name,
		&out.Currency,
		&out.Amount,
		&out.Uses,
		&out.Creator,
		&expiredAt,
		&createdAt)

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

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt.(time.Time)), timestamppb.New(createdAt.(time.Time))
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
		Type:     common.TransactionType_ActivatePromoCode,
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

	r.logger.Debugf("add in history: user: %d activate promo: %d", in.UserId, in.PromoId)
	return nil
}

// Delete activation of promo from table UserToPromo
func (r *Repository) DeleteActivatePromoFromHistory(
	ctx context.Context,
	tx pgx.Tx,
	in *promos.PromoId) (err error) {

	q := `DELETE FROM UserToPromo WHERE PromoId=$1`
	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
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

	if *activated {
		return true, Err.ErrPromoAlreadyActivated
	}

	return false, nil
}

// Check promo is expired or not
func (r *Repository) PromoIsExpired(in *promos.PromoCode) (b bool, err error) {
	if float64(time.Now().Unix()) > float64(in.ExpAt.AsTime().Unix()) {
		return false, Err.ErrPromoExpired
	}

	return true, nil
}

// Check promo in stock or not
func (r *Repository) PromoIsNotInStock(in *promos.PromoCode) (b bool, err error) {
	if in.Uses == 0 {
		return false, Err.ErrPromoNotInStock
	}

	return true, nil
}

// Check creator is owner or not
func (r *Repository) CreatorIsOwner(ctx context.Context, in *promos.PromoCode) (b bool, err error) {
	user, err := r.UserService.GetUser(ctx, &users.Id{Id: in.Creator})
	if err != nil {
		return false, Err.ErrServiceUsers
	}

	if user.Role == 3 {
		return true, nil
	}

	return false, Err.ErrCreatorIsNotOwner
}

// Add time for promo
func (r *Repository) AddTime(ctx context.Context, tx pgx.Tx, in *promos.AddTimeIn) (err error) {

	q := `UPDATE Promos SET ExpAt = $1 WHERE Id = $2`
	expAt := in.ExpAt.AsTime().Format("2006-01-02 15:04:05")

	if _, err := tx.Exec(ctx, q, expAt, in.PromoId); err != nil {
		return Err.ErrExecQuery
	}

	return nil
}

// Add uses for promo
func (r *Repository) AddUses(ctx context.Context, tx pgx.Tx, in *promos.AddUsesIn) (err error) {

	q := `UPDATE Promos SET Uses = Uses+$1 WHERE Id = $2`

	if _, err := tx.Exec(ctx, q, in.Uses, in.PromoId); err != nil {
		return Err.ErrExecQuery
	}

	return nil
}
