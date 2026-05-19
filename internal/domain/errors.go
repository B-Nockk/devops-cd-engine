package domain

import "errors"

var (
	ErrInvalidID       = errors.New("invalid ID")
	ErrInvalidStrategy = errors.New("invalid deployment strategy")
	ErrInvalidStatus   = errors.New("invalid status transition")
	ErrValidation      = errors.New("validation error")
)
