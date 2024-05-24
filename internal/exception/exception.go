package exception

import "errors"

var (
	ErrUnsupportedGrantType = errors.New("unsupported grant type")
	ErrPasswordIncorrect    = errors.New("password incorrect")

	ErrClientNotFound = errors.New("client not found")

	ErrUserNotFound = errors.New("user not found")

	ErrTokenNotFound         = errors.New("token not found")
	ErrTokenExpired          = errors.New("token is expired")
	ErrTokenUsedBeforeIssued = errors.New("token used before issued")
)
