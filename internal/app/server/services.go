package server

import (
	"github.com/alnovi/sso/internal/service"
	"github.com/alnovi/sso/internal/service/auth"
	"github.com/alnovi/sso/internal/service/client"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/internal/service/user"
)

type Services struct {
	auth   service.Auth
	client service.Client
	token  service.Token
	user   service.User
}

func newServices(app *App, a *Adapters) (*Services, error) {
	return &Services{
		auth:   auth.New(a.repo),
		client: client.New(a.repo),
		token:  token.New(a.repo),
		user:   user.New(a.repo),
	}, nil
}
