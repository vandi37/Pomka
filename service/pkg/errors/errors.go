package errors

import (
	"errors"
)

var (
	ErrUniquePromo         = errors.New("error name of promo must be unique")
	ErrExecQuery           = errors.New("error exeсution query")
	ErrTransactionCommit   = errors.New("error transaction commit")
	ErrTransactionRollback = errors.New("error transaction rollback")
)
