package database

import "errors"

const (
	ErrCodeInvalidConstraint = 1452
)

var (
	ErrNoRows          = errors.New("no rows in result set")
	ErrInvalidArgument = errors.New("invalid argument")
)
