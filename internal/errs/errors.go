package errs

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found from id")
	ErrUserIsUnableToPay = errors.New("user is unable to pay")
)
