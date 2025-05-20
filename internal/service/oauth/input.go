package oauth

type InputAuthorizeParams struct {
	ClientId     string
	ResponseType string
	RedirectUri  string
}

type InputAuthorizeBySession struct {
	ClientId     string
	ResponseType string
	RedirectUri  string
	State        string
	SessionId    string
}

type InputAuthorizeByCode struct {
	ClientId     string
	ResponseType string
	RedirectUri  string
	State        string
	Login        string
	Password     string
	UserIP       string
	UserAgent    string
}

type InputTokenByCode struct {
	ClientId     string
	ClientSecret string
	Code         string
}

type InputTokenByRefresh struct {
	ClientId     string
	ClientSecret string
	Refresh      string
}

type InputForgotPassword struct {
	ClientId    string
	RedirectUri string
	Query       string
	Login       string
	IP          string
	Agent       string
}

type InputResetPassword struct {
	Hash     string
	Password string
}
