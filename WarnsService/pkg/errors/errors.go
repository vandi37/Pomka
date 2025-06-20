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
	ErrUserIsNotModerator  = errors.New("error only moderators can give warns")
	ErrCreateWarn          = errors.New("error create warn")
	ErrCreateBan           = errors.New("error create ban")
	ErrGetWarns            = errors.New("error get warns")
	ErrGetBans             = errors.New("error get bans")
	ErrMakeWarnsInActive   = errors.New("error make warns inactive")
	ErrMakeBansInActive    = errors.New("error make bans inactive")
	ErrCountActiveWarns    = errors.New("error get count of active warns")
	ErrSendTransaction     = errors.New("error send transaction to service users")
	ErrUserAlreadyBanned   = errors.New("error user already banned")
)
