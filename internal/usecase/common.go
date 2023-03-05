package usecase

import (
	"github.com/c0depwn/fhnw-vesys-bank-server/internal/domain"
)

var (
	ErrOverdraw = Error{code: "overdraw"}
	ErrInactive = Error{code: "inactive"}
	ErrInvalid  = Error{code: "invalid"}
)

type Error struct {
	code  string
	msg   string
	cause error
}

func (e *Error) addCause(err error) {
	e.cause = err
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Error() string {
	return e.code
}

type DB interface {
	Get() (*domain.Bank, error)
}
