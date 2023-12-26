package dto

import "github.com/alnovi/sso/internal/entity"

const (
	GrantTypeAuthorizationCode = "authorization_code"
	GrantTypeRefreshToken      = "refresh_token"
)

type AccessToken struct {
	Client    entity.Client
	Code      string
	Refresh   string
	GrantType string
}
