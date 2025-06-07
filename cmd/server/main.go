package main

import (
	"github.com/alnovi/sso/internal/app/server"
)

// @title        SSO
// @description  Single sign-on (сервис единого входа)
// @version      1.0.0
//
// @contact.name Alnovi
// @contact.url  https://github.com/alnovi
//
// @license.name MIT
// @license.url  https://raw.githubusercontent.com/alnovi/sso/refs/heads/master/LICENSE.md
//
// @query.collection.format multi
func main() {
	server.NewApp(nil).Start(nil)
}
