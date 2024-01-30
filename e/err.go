package e

import "errors"

var (
	UserExists = errors.New("user already exists")

// ErrNotExists    = errors.New("row not exists")
// ErrUpdateFailed = errors.New("update failed")
// ErrDeleteFailed = errors.New("delete failed")
)
