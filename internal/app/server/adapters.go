package server

import (
	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/adapter/repository/postgres"
)

type Adapters struct {
	repo repository.Repository
}

func newAdapters(app *App) (*Adapters, error) {
	repo, err := postgres.NewRepository(postgres.Config{
		Host:     app.cfg.Database.Host,
		Port:     app.cfg.Database.Port,
		Database: app.cfg.Database.Database,
		User:     app.cfg.Database.User,
		Password: app.cfg.Database.Password,
		SSLMode:  app.cfg.Database.SSLMode,
	})

	return &Adapters{repo: repo}, err
}
