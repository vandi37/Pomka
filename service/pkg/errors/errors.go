package errors

import (
	"errors"
)

var (
	ErrExpAt                 = errors.New("error timestamb of expired data must be more then now")
	ErrValueUses             = errors.New("error value of uses must be > 0, for infinity uses set -1")
	ErrUniquePromo           = errors.New("error name of promo must be unique")
	ErrExecQuery             = errors.New("error exeсution query")
	ErrMissingPromoId        = errors.New("error missing promo id")
	ErrMissingPromoName      = errors.New("error missing promo name")
	ErrTransactionCommit     = errors.New("error transaction commit")
	ErrTransactionRollback   = errors.New("error transaction rollback")
	ErrMissingEnviroment     = errors.New("error missing enviroment")
	ErrWrongUserId           = errors.New("error wrong user id")
	ErrIncorrectData         = errors.New("error cannot scan data")
	ErrServiceUsers          = errors.New("error on service users")
	ErrWrongTypeData         = errors.New("error not supported type data")
	ErrPromoExpired          = errors.New("error promocode expired")
	ErrPromoNotInStock       = errors.New("error promocode activations are over")
	ErrPromoAlreadyActivated = errors.New("error promo is already activated by user")
)
