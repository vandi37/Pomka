package errors

import (
	"errors"
)

var (
	ErrUniquePromo         = errors.New("error name of promo must be unique")
	ErrExecQuery           = errors.New("error exeсution query")
	ErrTransactionCommit   = errors.New("error transaction commit")
	ErrTransactionRollback = errors.New("error transaction rollback")
	ErrMissingEnviroment   = errors.New("error missing enviroment")
	ErrWrongUserId         = errors.New("error wrong user id")
	ErrIncorrectData       = errors.New("error cannot scan data")
	ErrServiceUsers        = errors.New("error on service users")
	ErrWrongTypeData       = errors.New("error not supported type data")
)
