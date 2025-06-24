package mock

import (
	"checks/pkg/models/users"
	repeatible "checks/pkg/utils"
	"context"
	"fmt"

	Err "checks/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockServiceUsers struct {
	db *pgxpool.Pool
}

func NewMockServiceUsers(pool *pgxpool.Pool) *MockServiceUsers {
	return &MockServiceUsers{db: pool}
}

func (m *MockServiceUsers) GetUser(ctx context.Context, in *users.Id, opts ...grpc.CallOption) (*users.User, error) {
	var user = new(users.User)
	if errTx := repeatible.RunInTx(m.db, ctx, func(tx pgx.Tx) error {
		q := `SELECT * FROM Users WHERE Id=$1`
		if err := tx.QueryRow(ctx, q, in.Id).Scan(nil, nil, nil, nil, &user.Role, nil, nil, nil); err != nil {
			return err
		}

		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return user, nil
}

func (m *MockServiceUsers) SendTransaction(ctx context.Context, in *users.TransactionRequest, opts ...grpc.CallOption) (*users.TransactionResponse, error) {
	return nil, nil
}

func (m *MockServiceUsers) Create(ctx context.Context, role int, money int64) (int64, error) {
	var userId = new(int64)
	if errTx := repeatible.RunInTx(m.db, ctx, func(tx pgx.Tx) error {
		q := `INSERT INTO Users (Credits, Stocks, PremiumCredits, Role, AutoBuyEnabled, LastFarmingAt, CreatedAt) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING Id`
		if err := tx.QueryRow(ctx, q, money, 0, 0, role, false, timestamppb.Now().AsTime().Format("2006-01-02 15:04:05"), timestamppb.Now().AsTime().Format("2006-01-02 15:04:05")).Scan(&userId); err != nil {
			return Err.ErrExecQuery
		}

		return nil
	}); errTx != nil {
		return 0, errTx
	}

	return *userId, nil
}

func (m *MockServiceUsers) Delete(ctx context.Context, userId int64) (err error) {
	if errTx := repeatible.RunInTx(m.db, ctx, func(tx pgx.Tx) error {
		q := `DELETE FROM Users WHERE Id = $1`
		if _, err := tx.Exec(ctx, q, userId); err != nil {
			fmt.Println(err)
			return Err.ErrExecQuery
		}

		return nil
	}); errTx != nil {
		return errTx
	}

	return nil
}
