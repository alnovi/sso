package dto

import "github.com/alnovi/sso/internal/entity"

const (
	ResponseTypeCode = "code"
	GrantTypeCode    = "authorization_code"
	GrantTypeRefresh = "refresh_token"
)

type ValidateResponseType struct {
	ClientID     string
	ResponseType string
}

type AuthByCredentials struct {
	Client   *entity.Client
	Email    string
	Password string
}

type AccessTokenByCode struct {
	ClientID     string
	ClientSecret string
	CodeHash     string
	IP           string
	Agent        string
}

type AccessTokenByRefresh struct {
	ClientID     string
	ClientSecret string
	RefreshHash  string
	IP           string
	Agent        string
}

type ForgotPassword struct {
	Client *entity.Client
	Email  string
	IP     string
	Agent  string
}

type ResetPassword struct {
	Hash     string
	Password string
}
