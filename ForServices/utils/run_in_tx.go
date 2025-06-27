package utils

import (
	"context"

	e "errorspomka"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunInTx(db *pgxpool.Pool, ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	errFn := fn(tx)
	if errFn == nil {
		if errCommit := tx.Commit(ctx); errCommit != nil {
			return e.ErrTransactionCommit
		}
		return nil
	}

	if errRollback := tx.Rollback(ctx); errRollback != nil {
		return e.ErrTransactionRollback
	}

	return errFn
}
