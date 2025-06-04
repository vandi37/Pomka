package postgres

import (
	"context"
	"fmt"
	"time"

	repeatible "promos/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg DBConfig) (pool *pgxpool.Pool, err error) {
	time.Sleep(time.Duration(cfg.DelayAtmpsS))

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	err = repeatible.DoWithTries(func() error {

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, connStr)
		if err != nil {
			return fmt.Errorf("postgres: NewClient: error failed to connect to postgres")
		}

		return nil
	}, cfg.MaxAtmps, time.Duration(cfg.DelayAtmpsS)*time.Second)

	if err != nil {
		return nil, fmt.Errorf("postgres: NewClient: %s", err)
	}

	return pool, nil
}
