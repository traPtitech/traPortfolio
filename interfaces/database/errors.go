package database

import "errors"

var (
	ErrNoRows          = errors.New("no rows in result set")
	ErrInvalidArgument = errors.New("invalid argument")
)
