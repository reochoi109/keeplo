package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveAccount    = errors.New("account is inactive")
	ErrPasswordMismatch   = errors.New("password mismatch")
	ErrNicknameRequired   = errors.New("nickname is required")
	ErrAlreadyDeleted     = errors.New("user already deleted")
)
