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

	ErrNewPasswordTooWeak = errors.New("new password is too weak")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrUpdateFailed       = errors.New("user update failed")
	ErrDeleteFailed       = errors.New("user delete failed")
	ErrDatabase           = errors.New("database error")
	ErrInvalidInput       = errors.New("invalid input value")
	ErrInvalidEmailFormat = errors.New("invalid email format")
)
