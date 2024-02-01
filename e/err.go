package e

import "errors"

var (
	UserExists         = errors.New("user already exists")
	InvalidCredentials = errors.New("invalid credentials")

// ErrNotExists    = errors.New("row not exists")
// ErrUpdateFailed = errors.New("update failed")
// ErrDeleteFailed = errors.New("delete failed")
)
