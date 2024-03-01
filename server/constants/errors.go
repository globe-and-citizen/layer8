package constants

import "errors"

var (
	ErrUnknownRepositoryFactory = errors.New("unknown repository factory")
	ErrRequiredFieldMissing     = errors.New("required field missing")
	ErrInvalidField             = errors.New("invalid field")
	ErrInvalidEmail             = errors.New("invalid email")
	ErrAlreadyExists            = errors.New("already exists")
	ErrNotFound                 = errors.New("not found")
	ErrInvalidPassword          = errors.New("password is invalid")
	ErrMissingFields            = errors.New("missing fields")
)
