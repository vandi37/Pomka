package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg Config) (pool *pgxpool.Pool, err error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	// Trying several times
	err = utils.DoWithTries(func() error {

		// After 5 seconds, close ctx
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// Connecting to postgres pool
		pool, err = pgxpool.New(ctx, connStr)
		if err != nil {
			return fmt.Errorf("error failed connect to postgres pool")
		}

		return nil
	}, cfg.MaxAtmps, time.Duration(cfg.DelayAtmpsS)*time.Second)

	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("cannot connecting to postgres pool with `postgres://%s:%s@%s:%s/%s`",
			cfg.User, "<PASSWORD_SECRET>", cfg.Host, cfg.Port, cfg.Database))
	}

	return pool, nil
}
