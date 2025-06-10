package mock

import (
	"context"
	"fmt"
	"promos/internal/models/users"
	repeatible "promos/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockServiceUsers struct {
	db *pgxpool.Pool
	mock.Mock
}

func NewMockServiceUsers(pool *pgxpool.Pool) *MockServiceUsers {
	return &MockServiceUsers{db: pool}
}

func (m *MockServiceUsers) SendTransaction(ctx context.Context, in *users.TransactionRequest, opts ...grpc.CallOption) (*users.TransactionResponse, error) {
	return nil, nil
}

func (m *MockServiceUsers) Create(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockServiceUsers) Delete(ctx context.Context, userId int64) (err error) {
	if errTx := repeatible.RunInTx(m.db, ctx, func(tx pgx.Tx) error {
		q := `DELETE FROM Users WHERE Id = $1`
		if _, err := tx.Exec(ctx, q, userId); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	}); errTx != nil {
		return errTx
	}

	return nil
}
