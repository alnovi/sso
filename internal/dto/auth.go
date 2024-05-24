package dto

import "github.com/alnovi/sso/internal/entity"

const (
	ResponseTypeCode = "code"
	GrantTypeCode    = "authorization_code"
	GrantTypeRefresh = "refresh_token"
)

type InpValidateGrantType struct {
	ClientID    string
	GrantType   string
	RedirectURI string
}

type InpAuthByCredentials struct {
	Client   *entity.Client
	Email    string
	Password string
}
