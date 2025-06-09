package repository

import (
	"context"
	"promos/internal/models/common"
	"promos/internal/models/promos"
	"promos/internal/models/users"
	"promos/internal/transport/grpc/conn"
	Err "promos/pkg/errors"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	*conn.ClientsServices
}

func NewRepository(conn *conn.ClientsServices) *Repository {
	return &Repository{conn}
}

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

func (r *Repository) Delete(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (*common.Response, error) {
	q := `DELETE FROM Promos WHERE Id = $1`

	if _, err := tx.Exec(ctx, q, in.Id); err != nil {
		return nil, Err.ErrExecQuery
	}

	return nil, nil
}

func (r *Repository) Use(ctx context.Context, tx pgx.Tx, in *promos.PromoUserId) (*users.TransactionResponse, error) {
	promo, err := r.GetPromoById(ctx, tx, &promos.PromoId{Id: in.PromoId})
	if err != nil {
		return nil, err
	}

	if _, err := r.ClientsServices.UsersClient.SendTransaction(ctx, &users.TransactionRequest{
		Sender:   nil,
		Receiver: &users.UserTransaction{UserId: in.UserId, Amount: promo.Amount, Currency: promo.Currency},
	}); err != nil {
		return nil, Err.ErrServiceUsers
	}

	return nil, nil
}

func (r *Repository) GetPromoById(ctx context.Context, tx pgx.Tx, in *promos.PromoId) (*promos.PromoCode, error) {
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Id = $1`
	row := tx.QueryRow(ctx, q, in.Id)

	if err := row.Scan(&out.Id, &out.Name, &out.Currency, &out.Amount, &out.Uses, &out.Creator, nil, nil); err != nil {
		return nil, Err.ErrIncorrectData
	}
	return out, nil
}

func (r *Repository) GetPromoByName(ctx context.Context, tx pgx.Tx, in *promos.PromoName) (*promos.PromoCode, error) {
	var out = new(promos.PromoCode)

	q := `SELECT * FROM Promos WHERE Name = $1`
	row := tx.QueryRow(ctx, q, in.Name)

	if err := row.Scan(&out.Id, &out.Name, &out.Currency, &out.Amount, &out.Uses, &out.Creator, nil, nil); err != nil {
		return nil, Err.ErrIncorrectData
	}

	return out, nil
}
