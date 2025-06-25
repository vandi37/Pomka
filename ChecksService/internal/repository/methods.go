package repository

import (
	Err "checks/pkg/errors"
	"checks/pkg/models/checks"
	"checks/pkg/models/users"
	"checks/pkg/postgres"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Repository) CreateCheck(ctx context.Context, db postgres.DB, in *checks.CheckCreate) (*checks.Check, error) {
	var check = new(checks.Check)
	var createdAt = new(time.Time)
	var key = uuid.New().String()

	q := `INSERT INTO "Checks" ("CreatorId", "Key", "Currency", "Amount") 
		  VALUES ($1, $2, $3, $4)
		  RETURNING "Id", "CreatorId", "Key", "Currency", "Amount", "CreatedAt"`

	if err := db.QueryRow(
		ctx, q, in.Creator, r.h.Hash(key), in.Currency, in.Amount).
		Scan(&check.Id, &check.Creator, nil, &check.Currency, &check.Amount, &createdAt); err != nil {
		return nil, errors.Join(Err.ErrExecQuery, err)
	}

	check.CreatedAt = timestamppb.New(*createdAt)
	check.Key = key

	return check, nil
}

func (r *Repository) RemoveCheck(ctx context.Context, db postgres.DB, in *checks.CheckId) error {
	q := `DELETE FROM "Checks"
	      WHERE "Id"=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(Err.ErrExecQuery, err)
	}

	return nil
}

func (r *Repository) GetUsersCheck(ctx context.Context, db postgres.DB, in *users.Id) (*checks.AllChecks, error) {
	var allChecks = new(checks.AllChecks)

	q := `SELECT * FROM "Checks"
	      WHERE "CreatorId"=$1`

	rows, err := db.Query(ctx, q, in.Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, errors.Join(Err.ErrExecQuery, err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var check = new(checks.Check)
		var createdAt = new(time.Time)

		if err := rows.Scan(&check.Id, &check.Creator, &check.Key, &check.Currency, &check.Amount, &createdAt); err != nil {
			return nil, errors.Join(Err.ErrIncorrectData, err)
		}

		check.CreatedAt = timestamppb.New(*createdAt)
		allChecks.Checks = append(allChecks.Checks, check)
	}

	return allChecks, nil
}

func (r *Repository) GetCheckByKey(ctx context.Context, db postgres.DB, key string) (*checks.Check, error) {
	var check = new(checks.Check)
	var createdAt = new(time.Time)
	key = r.h.Hash(key)

	q := `SELECT * FROM "Checks"
	      WHERE "Key"=$1`

	err := db.QueryRow(
		ctx, q, key).
		Scan(&check.Id, &check.Creator, &check.Key, &check.Currency, &check.Amount, createdAt)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.Join(Err.ErrCheckNotValid, err)
		default:
			return nil, errors.Join(Err.ErrExecQuery, err)
		}
	}

	check.CreatedAt = timestamppb.New(*createdAt)
	return check, nil
}
