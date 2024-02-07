package e

import "errors"

var (
	ErrSystem             = errors.New("system error")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

// ErrNotExists    = errors.New("row not exists")
// ErrUpdateFailed = errors.New("update failed")
// ErrDeleteFailed = errors.New("delete failed")
)
