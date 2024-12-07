package errors

import "errors"

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrInvalidGUID               = errors.New("invalid guid")
	ErrUserAlreadyHasAccessToken = errors.New("user already has a access token")
	ErrUnauthorized              = errors.New("no user with this token")
	ErrCreateRefreshToken        = errors.New("can't create refresh token")
	ErrNoArguments               = errors.New("no arguments")
	ErrOnlyOneArgument           = errors.New("must be only one argument")
	ErrRefreshTokenNotUnique     = errors.New("refresh token is not unique")
	ErrSendWarning               = errors.New("can't send warning")
)
