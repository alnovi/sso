package main

import (
	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/docs"
	"github.com/alnovi/sso/internal/app/server"
	"github.com/alnovi/sso/pkg/configure"
	"github.com/alnovi/sso/pkg/logger"
)

// @title          Single Sign-On Server
// @version        0.0.0
// @description    Открытый протокол, обеспечивающий безопасную авторизацию простым и стандартным способом.
//
// @license.name MIT License
// @license.url https://mit-license.org/
//
// @host 127.0.0.1:8081
// @BasePath /
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @query.collection.format multi
// @schemes http https
func main() {
	var err error

	cfg, err := configure.LoadFromEnv[config.Config](config.NewConfig())
	must(err)

	log := logger.NewLogger(cfg.Logger.Format, cfg.Logger.Level)
	docs.SwaggerInfo.Host = cfg.App.Host

	app, err := server.NewApp(&cfg, log)
	must(err)

	err = app.Run()
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
