package main

import (
	"github.com/alnovi/sso/internal/app/server"
)

// @title        SSO
// @description  Single sign-on (сервис единого входа)
// @host         http://localhost:8081
// @base.path    /
// @version      1.0.0
//
// @contact.name Alnovi
// @contact.url  https://github.com/alnovi
//
// @license.name MIT
// @license.url  https://raw.githubusercontent.com/alnovi/sso/refs/heads/master/LICENSE.md
//
// @securitydefinitions.apikey JWT-Access
// @in header
// @name Authorization
//
// @query.collection.format multi
func main() {
	server.NewApp(nil).Start(nil)
}
