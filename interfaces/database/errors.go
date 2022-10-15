package database

import "errors"

var (
	ErrCodeInvalidConstraint = uint16(1452)
	ErrNoRows                = errors.New("no rows in result set")
	ErrInvalidArgument       = errors.New("invalid argument")
)
