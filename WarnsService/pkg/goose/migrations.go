package migrations

import (
	"context"
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Up(ctx context.Context, pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

func Down(ctx context.Context, pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.DownContext(ctx, db, "migrations"); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
