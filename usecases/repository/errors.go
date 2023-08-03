package repository

import (
	"errors"
)

var (
	// ErrNilID id is nil
	ErrNilID = errors.New("nil id")
	// ErrInvalidID id is invalid
	ErrInvalidID = errors.New("invalid uuid")
	// ErrNotFound not found
	ErrNotFound = errors.New("not found")
	// ErrForbidden forbidden
	ErrForbidden = errors.New("forbidden")
	// ErrAlreadyExists already exists
	ErrAlreadyExists = errors.New("already exists")
	// ErrInvalidArg argument error
	ErrInvalidArg = errors.New("argument error")
	// ErrDBInternal database internal error
	ErrDBInternal = errors.New("database internal error")
	// ErrBind failed to bind request
	ErrBind = errors.New("bind error")
	// ErrValidate failed to validate request
	ErrValidate = errors.New("validate error")
)
