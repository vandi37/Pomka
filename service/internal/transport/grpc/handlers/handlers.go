package service

import (
	"context"
	"warns/pkg/models/common"
	"warns/pkg/models/users"
	"warns/pkg/models/warns"
	repeatible "warns/pkg/utils"

	"github.com/jackc/pgx/v5"
)

// insert in table Warns
func (sw *ServiceWarns) Create(ctx context.Context, in *warns.CreateIn) (*warns.CreateFailure, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

// update warns for this user IsActive=false, insert in table Bans, query to service users for update role -> banned
func (sw *ServiceWarns) Ban(ctx context.Context, in *users.Id) (*warns.BanFailure, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}

// update ban for this user IsActive=false, query to service users for update role -> user
func (sw *ServiceWarns) Unban(ctx context.Context, in *users.Id) (*common.Response, error) {
	// Run in transaction
	if errTx := repeatible.RunInTx(sw.db, ctx, func(tx pgx.Tx) error {
		return nil
	}); errTx != nil {
		return nil, errTx
	}

	return nil, nil
}
