package repeatible

import (
	"context"
	"time"

	Err "promos/pkg/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DoWithTries(fn func() error, attemps int, delay time.Duration) (err error) {
	for attemps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemps--

			continue
		}

		return nil
	}

	return
}

func RunInTx(db *pgxpool.Pool, ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	errFn := fn(tx)
	if errFn == nil {
		if errCommit := tx.Commit(ctx); errCommit != nil {
			return Err.ErrTransactionCommit
		}
		return nil
	}

	if errRollback := tx.Rollback(ctx); errRollback != nil {
		return Err.ErrTransactionRollback
	}

	return errFn
}
