package repository

import (
	"context"
	"errors"
	Err "promos/pkg/errors"
	"promos/pkg/models/promos"
	"promos/pkg/models/users"
	"promos/pkg/postgres"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Insert promo to table promos
func (r *Repository) CreatePromo(
	ctx context.Context,
	db postgres.DB,
	in *promos.CreatePromo) (*promos.PromoCode, error) {

	// Check args valid, because Exec() send panic if have error
	if !(0 <= in.Currency && in.Currency <= 2 && in.Amount >= 0 && (in.Uses > 0 || in.Uses == -1)) {
		return nil, Err.ErrBadArgs
	}

	var expiredAt, createdAt time.Time // Scan() cannot convert sql timestamp to protobuf/types/known/timestamppb
	var out = new(promos.PromoCode)

	q := `INSERT INTO Promos (Name, Currency, Amount, Uses, Creator, ExpAt)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING Id, Name, Currency, Amount, Uses, Creator, ExpAt, CreatedAt`
	expAt := in.ExpAt.AsTime().Format("2006-01-02 15:04:05")

	err := db.QueryRow(ctx, q,
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
		return nil, errors.Join(Err.ErrExecQuery, err)
	}

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt), timestamppb.New(createdAt)
	return out, nil
}

// Delete promo from table promos by ID
func (r *Repository) DeletePromoById(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoId) error {
	q := `DELETE FROM Promos WHERE Id = $1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}

// Delete promo from table promos by Name
func (r *Repository) DeletePromoByName(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoName) error {
	q := `DELETE FROM Promos WHERE Name = $1`

	if _, err := db.Exec(ctx, q, in.Name); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}

// Get promo from table promos by ID
func (r *Repository) GetPromoById(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoId) (*promos.PromoCode, error) {

	var expiredAt, createdAt time.Time
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Id = $1`
	row := db.QueryRow(ctx, q, in.Id)

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
			return nil, errors.Join(Err.ErrMissingPromoId, err)
		default:
			return nil, errors.Join(Err.ErrExecQuery)
		}
	}

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt), timestamppb.New(createdAt)
	return out, nil
}

// Get promo from table promos by Name
func (r *Repository) GetPromoByName(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoName) (*promos.PromoCode, error) {

	var expiredAt, createdAt time.Time
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Name = $1`
	row := db.QueryRow(ctx, q, in.Name)

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
			return nil, errors.Join(Err.ErrMissingPromoId, err)
		default:
			return nil, errors.Join(Err.ErrExecQuery)
		}
	}

	out.ExpAt, out.CreatedAt = timestamppb.New(expiredAt), timestamppb.New(createdAt)
	return out, nil
}

// Update table promos, decrement uses of promo.
func (r *Repository) DecrementPromoUses(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoId) (err error) {
	if in.Id == -1 {
		return nil
	}

	q := `UPDATE Promos SET Uses = Uses-1 WHERE Id = $1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}

// Insert activation of promo to table UserToPromo
func (r *Repository) AddActivatePromoToHistory(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoUserId) (err error) {
	q := `INSERT INTO UserToPromo (UserId, PromoId, ActivatedAt) VALUES ($1, $2, $3)`

	actAt := time.Now().Format("2006-01-02 15:04:05")
	if _, err := db.Exec(ctx, q, in.UserId, in.PromoId, actAt); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}
	return nil
}

// Delete activation of promo from table UserToPromo
func (r *Repository) DeleteActivatePromoFromHistory(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoId) (err error) {

	q := `DELETE FROM UserToPromo WHERE PromoId=$1`
	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}

// If promo been activated by user, return true. If promo not activated by user, return false.
func (r *Repository) PromoIsAlreadyActivated(
	ctx context.Context,
	db postgres.DB,
	in *promos.PromoUserId) (b bool, err error) {

	var activated = new(bool)

	q := `SELECT EXISTS(SELECT * FROM UserToPromo WHERE UserId = $1 AND PromoId = $2)`

	row := db.QueryRow(ctx, q, in.UserId, in.PromoId)
	if err := row.Scan(&activated); err != nil {
		return false, errors.Join(Err.ErrExecQuery, err)
	}

	if *activated {
		return true, errors.Join(Err.ErrPromoAlreadyActivated, err)
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
func (r *Repository) CreatorIsOwner(ctx context.Context, user *users.User) (b bool, err error) {

	if user.Role == 3 {
		return true, nil
	}

	return false, Err.ErrCreatorIsNotOwner
}

// Add time for promo
func (r *Repository) AddTime(ctx context.Context, db postgres.DB, in *promos.AddTimeIn) (err error) {

	q := `UPDATE Promos SET ExpAt = $1 WHERE Id = $2`
	expAt := in.ExpAt.AsTime().Format("2006-01-02 15:04:05")

	if _, err := db.Exec(ctx, q, expAt, in.PromoId); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}

// Add uses for promo
func (r *Repository) AddUses(ctx context.Context, db postgres.DB, in *promos.AddUsesIn) (err error) {

	q := `UPDATE Promos SET Uses = Uses+$1 WHERE Id = $2`

	if _, err := db.Exec(ctx, q, in.Uses, in.PromoId); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}
