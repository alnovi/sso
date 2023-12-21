package exception

var (
	PasswordIncorrect = New("password incorrect")
	NotAuthorization  = New("not authorization")

	ClientNotFound     = New("client not found")
	ClientAccessDenied = New("access denied")

	TokenNotFound = New("token not found")

	UserNotFound = New("user not found")
)
