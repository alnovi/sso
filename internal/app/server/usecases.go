package server

import (
	"github.com/alnovi/sso/internal/usecase"
	"github.com/alnovi/sso/internal/usecase/auth"
	"github.com/alnovi/sso/internal/usecase/client"
)

type UseCases struct {
	auth   usecase.Auth
	client usecase.Client
}

func newUseCases(_ *App, s *Services) (*UseCases, error) {
	return &UseCases{
		auth:   auth.New(s.auth, s.token, s.user),
		client: client.New(s.client),
	}, nil
}
