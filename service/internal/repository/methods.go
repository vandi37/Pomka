package repository

import (
	"context"
	"time"
	Err "warns/pkg/errors"
	"warns/pkg/models/users"
	"warns/pkg/models/warns"
	"warns/pkg/postgres"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *Repository) CreateWarn(ctx context.Context, db postgres.DB, in *warns.WarnCreate) (warn *warns.Warn, err error) {
	warn = new(warns.Warn)
	var issuedAt interface{}

	q := `INSERT INTO Warns (UserId, ModeratorId, Reason) 
		  VALUES ($1, $2, $3)
		  RETURNING Id, UserId, ModeratorId, Reason, IssuedAt, IsActive`

	err = db.QueryRow(ctx, q, in.UserId, in.ModerId, in.Reason).Scan(&warn.Id, &warn.UserId, &warn.ModerId, &warn.Reason, &issuedAt, &warn.IsActive)
	if err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return nil, Err.ErrExecQuery
	}

	warn.IssuedAt = timestamppb.New(issuedAt.(time.Time))

	return warn, nil
}

func (r *Repository) CreateBan(ctx context.Context, db postgres.DB, in *warns.BanCreate) (ban *warns.Ban, err error) {
	ban = new(warns.Ban)
	var issuedAt interface{}

	q := `INSERT INTO Bans (UserId, ModeratorId, Reason) 
		  VALUES ($1, $2, $3)
		  RETURNING Id, UserId, ModeratorId, Reason, IssuedAt, IsActive`

	err = db.QueryRow(ctx, q, in.UserId, in.ModerId, in.Reason).Scan(&ban.Id, &ban.UserId, &ban.ModerId, &ban.Reason, &issuedAt, &ban.IsActive)
	if err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return nil, Err.ErrExecQuery
	}

	ban.IssuedAt = timestamppb.New(issuedAt.(time.Time))

	return ban, nil
}

func (r *Repository) GetWarns(ctx context.Context, db postgres.DB, in *users.Id) (allwarns *warns.AllWarns, err error) {
	q := `SELECT * FROM Warns
		  WHERE UserId=$1`

	rows, err := db.Query(ctx, q, in.Id)
	if err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return nil, Err.ErrExecQuery
	}
	defer rows.Close()

	for rows.Next() {
		var warn = new(warns.Warn)
		var issuedAt interface{}

		if err := rows.Scan(&warn.Id, &warn.UserId, &warn.ModerId, &warn.Reason, &issuedAt, &warn.IsActive); err != nil {
			r.logger.Warn(Err.ErrIncorrectData, err)
			return nil, Err.ErrIncorrectData
		}

		warn.IssuedAt = timestamppb.New(issuedAt.(time.Time))

		allwarns.Warns = append(allwarns.Warns, warn)
	}

	return allwarns, nil
}

func (r *Repository) GetBans(ctx context.Context, db postgres.DB, in *users.Id) (allbans *warns.AllBans, err error) {
	q := `SELECT * FROM Bans
	  WHERE UserId=$1`

	rows, err := db.Query(ctx, q, in.Id)
	if err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return nil, Err.ErrExecQuery
	}
	defer rows.Close()

	for rows.Next() {
		var ban = new(warns.Ban)
		var issuedAt interface{}

		if err := rows.Scan(&ban.Id, &ban.UserId, &ban.ModerId, &ban.Reason, &issuedAt, &ban.IsActive); err != nil {
			r.logger.Warn(Err.ErrIncorrectData, err)
			return nil, Err.ErrIncorrectData
		}

		ban.IssuedAt = timestamppb.New(issuedAt.(time.Time))

		allbans.Bans = append(allbans.Bans, ban)
	}

	return allbans, nil
}

func (r *Repository) DeleteHistoryWarns(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `DELETE FROM Warns WHERE UserId=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return Err.ErrExecQuery
	}

	return nil
}

func (r *Repository) DeleteHistoryBans(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `DELETE FROM Warns WHERE UserId=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return Err.ErrExecQuery
	}

	return nil
}

func (r *Repository) IsUserModerator(ctx context.Context, in *users.User) (b bool, err error) {

	if in.Role == 2 {
		return true, nil
	}

	return false, Err.ErrUserIsNotModerator
}

func (r *Repository) MakeWarnsInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `UPDATE Warns SET IsActive=FALSE WHERE UserId=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return Err.ErrExecQuery
	}

	return nil
}

func (r *Repository) MakeBanInActive(ctx context.Context, db postgres.DB, in *users.Id) (err error) {
	q := `UPDATE Bans SET IsActive=FALSE WHERE UserId=$1`

	if _, err := db.Exec(ctx, q, in.Id); err != nil {
		r.logger.Warn(Err.ErrExecQuery, err)
		return Err.ErrExecQuery
	}

	return nil
}

func (r *Repository) GetCountOfActiveWarns(ctx context.Context, db postgres.DB, in *users.Id) (int, error) {
	var cnt = new(int)

	q := `SELECT COUNT(*) FROM Warns WHERE UserId=$1 AND IsActive=TRUE`
	if err := db.QueryRow(ctx, q, in.Id).Scan(&cnt); err != nil {
		return 0, Err.ErrExecQuery
	}

	return *cnt, nil
}
