package repository

import (
	"context"
	"errors"
	e "errorspomka"
	"postgres"
	"protobuf/users"
	"protobuf/warns"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Repository) CreateWarn(ctx context.Context, db postgres.DB, in *warns.ModerUserReason) (warn *warns.Warn, err error) {
	warn = new(warns.Warn)
	var issuedAt = new(time.Time)

	q := `INSERT INTO "Warns" ("UserId", "ModeratorId", "Reason") 
		  VALUES ($1, $2, $3)
		  RETURNING "Id", "UserId", "ModeratorId", "Reason", "IssuedAt", "IsActive"`

	err = db.QueryRow(ctx, q, in.UserId, in.ModerId, in.Reason).Scan(&warn.Id, &warn.UserId, &warn.ModerId, &warn.Reason, &issuedAt, &warn.IsActive)
	if err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}

	warn.IssuedAt = timestamppb.New(*issuedAt)

	return warn, nil
}

func (r *Repository) CreateBan(ctx context.Context, db postgres.DB, in *warns.ModerUserReason) (ban *warns.Ban, err error) {
	ban = new(warns.Ban)
	var issuedAt = new(time.Time)

	q := `INSERT INTO "Bans" ("UserId", "ModeratorId", "Reason") 
		  VALUES ($1, $2, $3)
		  RETURNING "Id", "UserId", "ModeratorId", "Reason", "IssuedAt", "IsActive"`

	err = db.QueryRow(ctx, q, in.UserId, in.ModerId, in.Reason).Scan(&ban.Id, &ban.UserId, &ban.ModerId, &ban.Reason, &issuedAt, &ban.IsActive)
	if err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}

	ban.IssuedAt = timestamppb.New(*issuedAt)

	return ban, nil
}

func (r *Repository) GetWarns(ctx context.Context, db postgres.DB, in *users.Id) (allwarns *warns.AllWarns, err error) {
	q := `SELECT * FROM "Warns"
		  WHERE "UserId"=$1`

	rows, err := db.Query(ctx, q, in.Id)
	if err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}
	defer rows.Close()

	for rows.Next() {
		var warn = new(warns.Warn)
		var issuedAt = new(time.Time)

		if err := rows.Scan(&warn.Id, &warn.UserId, &warn.ModerId, &warn.Reason, &issuedAt, &warn.IsActive); err != nil {
			return nil, errors.Join(e.ErrIncorrectData, err)
		}

		warn.IssuedAt = timestamppb.New(*issuedAt)

		allwarns.Warns = append(allwarns.Warns, warn)
	}

	return allwarns, nil
}

func (r *Repository) GetBans(ctx context.Context, db postgres.DB, in *users.Id) (allbans *warns.AllBans, err error) {
	q := `SELECT * FROM "Bans"
	  WHERE "UserId"=$1`

	rows, err := db.Query(ctx, q, in.Id)
	if err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}
	defer rows.Close()

	for rows.Next() {
		var ban = new(warns.Ban)
		var issuedAt = new(time.Time)

		if err := rows.Scan(&ban.Id, &ban.UserId, &ban.ModerId, &ban.Reason, &issuedAt, &ban.IsActive); err != nil {
			return nil, errors.Join(e.ErrIncorrectData, err)
		}

		ban.IssuedAt = timestamppb.New(*issuedAt)

		allbans.Bans = append(allbans.Bans, ban)
	}

	return allbans, nil
}

func (r *Repository) DeleteHistoryWarns(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `DELETE FROM "Warns"
		  WHERE "UserId"=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(e.ErrExecQuery, err)
	}

	return nil
}

func (r *Repository) DeleteHistoryBans(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `DELETE FROM "Bans"
		  WHERE "UserId"=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(e.ErrExecQuery, err)
	}

	return nil
}

func (r *Repository) IsUserModerator(ctx context.Context, in *users.User) (b bool, err error) {

	if in.Role == 2 {
		return true, nil
	}
	return false, errors.Join(e.ErrUserIsNotModerator, err)
}

func (r *Repository) MakeWarnsInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `UPDATE "Warns" 
          SET "IsActive"=FALSE
		  WHERE "UserId"=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(e.ErrExecQuery, err)
	}

	return nil
}

func (r *Repository) MakeBanInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `UPDATE "Bans" SET "IsActive"=FALSE WHERE "UserId"=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(e.ErrExecQuery, err)
	}

	return nil
}

func (r *Repository) GetCountOfActiveWarns(ctx context.Context, db postgres.DB, in *users.Id) (*warns.CountOfActiveWarns, error) {
	var cnt = new(int)

	q := `SELECT COUNT("Id") FROM "Warns" WHERE "UserId"=$1 AND "IsActive"=TRUE`
	if err := db.QueryRow(ctx, q, in.Id).Scan(&cnt); err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}

	return &warns.CountOfActiveWarns{CountWarns: int32(*cnt)}, nil
}

func (r *Repository) GetActiveWarns(ctx context.Context, db postgres.DB, in *users.Id) (allwarns *warns.AllWarns, err error) {
	q := `SELECT * FROM "Warns"
		  WHERE "UserId"=$1 AND "IsActive"=TRUE`

	rows, err := db.Query(ctx, q, in.Id)
	if err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}
	defer rows.Close()

	for rows.Next() {
		var warn = new(warns.Warn)
		var issuedAt = new(time.Time)

		if err := rows.Scan(&warn.Id, &warn.UserId, &warn.ModerId, &warn.Reason, &issuedAt, &warn.IsActive); err != nil {
			return nil, errors.Join(e.ErrIncorrectData, err)
		}

		warn.IssuedAt = timestamppb.New(*issuedAt)

		allwarns.Warns = append(allwarns.Warns, warn)
	}

	return allwarns, nil
}

func (r *Repository) GetActiveBan(ctx context.Context, db postgres.DB, in *users.Id) (ban *warns.Ban, err error) {
	q := `SELECT * FROM "Bans"
	      WHERE "UserId"=$1 AND "IsActive"=TRUE`

	var issuedAt = new(time.Time)
	ban = new(warns.Ban)

	if err := db.QueryRow(ctx, q, in.Id).Scan(&ban.Id, &ban.UserId, &ban.ModerId, &ban.Reason, &issuedAt, &ban.IsActive); err != nil {
		return nil, errors.Join(e.ErrExecQuery, err)
	}

	ban.IssuedAt = timestamppb.New(*issuedAt)

	return ban, nil
}

func (r *Repository) MakeLastWarnInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `DELETE FROM "Warns"
		  WHERE "UserId" = $1 AND "Id" = (SELECT "Id" FROM "Warns" ORDER BY "IssuedAt" DESC LIMIT 1);`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		return errors.Join(e.ErrExecQuery, err)
	}

	return nil
}

func (r *Repository) IsAlreadyBanned(ctx context.Context, db postgres.DB, in *users.Id) (bool, error) {
	var b = new(bool)

	q := `SELECT EXISTS(SELECT * FROM "Bans" WHERE "UserId"=$1 AND "IsActive"=TRUE)`

	if err := db.QueryRow(ctx, q, in.Id).Scan(&b); err != nil {
		return false, errors.Join(e.ErrExecQuery, err)
	}

	if *b {
		return true, e.ErrUserAlreadyBanned
	}

	return false, nil
}
