package errors

import "errors"

var (
	ErrInvalidDirectory = errors.New("invalid target directory")
	ErrInvalidURL       = errors.New("invalid repo url")
	ErrContextMismatch  = errors.New("missing or invalid type stored in context key")
)
