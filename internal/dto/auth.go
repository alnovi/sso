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
	RedirectURI  string
}

type AuthByCredentials struct {
	Client   *entity.Client
	Email    string
	Password string
}

type AccessTokenByGrantType struct {
	GrantType    string
	ClientID     string
	ClientSecret string
	TokenHash    string
}

type AccessTokenByCode struct {
	ClientID     string
	ClientSecret string
	CodeHash     string
}

type AccessTokenByRefresh struct {
	ClientID     string
	ClientSecret string
	RefreshHash  string
}
