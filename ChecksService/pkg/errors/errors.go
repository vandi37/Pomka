package errors

import (
	"errors"
)

var (
	ErrExecQuery           = errors.New("error exeсution query")
	ErrTransactionCommit   = errors.New("error transaction commit")
	ErrTransactionRollback = errors.New("error transaction rollback")
	ErrMissingEnviroment   = errors.New("error missing enviroment")
	ErrWrongUserId         = errors.New("error wrong user id")
	ErrIncorrectData       = errors.New("error cannot scan data")
	ErrServiceUsers        = errors.New("error on service users")
	ErrWrongTypeData       = errors.New("error not supported type data")
	ErrSendTransaction     = errors.New("error send transaction to service users")
	ErrCheckNotValid       = errors.New("check key invalid or missing")
)
