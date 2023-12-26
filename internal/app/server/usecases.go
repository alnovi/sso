package server

import (
	"github.com/alnovi/sso/internal/usecase"
	"github.com/alnovi/sso/internal/usecase/auth"
	"github.com/alnovi/sso/internal/usecase/client"
	"github.com/alnovi/sso/internal/usecase/token"
	"github.com/alnovi/sso/internal/usecase/user"
)

type UseCases struct {
	auth   usecase.Auth
	client usecase.Client
	token  usecase.Token
	user   usecase.User
}

func newUseCases(_ *App, s *Services) (*UseCases, error) {
	return &UseCases{
		auth:   auth.New(s.auth, s.token, s.user),
		client: client.New(s.client),
		token:  token.New(s.token, s.user),
		user:   user.New(s.user),
	}, nil
}
