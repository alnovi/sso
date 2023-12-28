package main

import (
	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/internal/app/server"
	"github.com/alnovi/sso/pkg/configure"
	"github.com/alnovi/sso/pkg/logger"
)

func main() {
	var err error

	cfg, err := configure.LoadFromEnv[config.Config](config.NewConfig())
	must(err)

	app, err := server.NewApp(&cfg, logger.NewLogger(cfg.Logger.Format, cfg.Logger.Level))
	must(err)

	err = app.Run()
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
