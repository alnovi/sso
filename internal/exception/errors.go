package exception

var (
	ErrPasswordIncorrect = New("password incorrect")
	ErrNotAuthorization  = New("not authorization")

	ErrClientNotFound     = New("client not found")
	ErrClientAccessDenied = New("access denied")

	ErrTokenNotFound = New("token not found")

	ErrUserNotFound = New("user not found")
)
