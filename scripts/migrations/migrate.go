package migrations

import "embed"

//go:embed *.sql
var Schema embed.FS
