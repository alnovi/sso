package exception

var (
	AccessDenied      = New("access denied")
	PasswordIncorrect = New("password incorrect")

	ClientNotFound = New("client not found")

	TokenNotFound = New("token not found")

	UserNotFound = New("user not found")
)
