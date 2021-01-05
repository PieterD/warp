package ecs

import (
	"errors"
	"fmt"
)

type userTemporary struct {
	cause error
}

func (e userTemporary) Error() string {
	return fmt.Sprintf("user error: %v", e.cause)
}

func NewUserError(cause error) error {
	if cause == nil {
		cause = fmt.Errorf("no cause specified")
	}
	return userTemporary{
		cause: cause,
	}
}

func IsUserError(err error) (cause error, ok bool) {
	var e userTemporary
	if !errors.As(err, &e) {
		return nil, false
	}
	return e.cause, true
}
