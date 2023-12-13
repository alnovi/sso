package dto

import "github.com/alnovi/sso/internal/entity"

type AccessToken struct {
	Client    entity.Client
	Code      string
	Refresh   string
	GrantType string
}
