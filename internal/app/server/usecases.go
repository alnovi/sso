package server

import (
	"github.com/alnovi/sso/internal/usecase"
	"github.com/alnovi/sso/internal/usecase/auth"
	"github.com/alnovi/sso/internal/usecase/client"
	"github.com/alnovi/sso/internal/usecase/profile"
	"github.com/alnovi/sso/internal/usecase/token"
)

type UseCases struct {
	auth    usecase.Auth
	client  usecase.Client
	token   usecase.Token
	profile usecase.Profile
}

func newUseCases(_ *App, s *Services) (*UseCases, error) {
	return &UseCases{
		auth:    auth.New(s.auth, s.token, s.user),
		client:  client.New(s.client),
		token:   token.New(s.token, s.user),
		profile: profile.New(s.client, s.token, s.user),
	}, nil
}
