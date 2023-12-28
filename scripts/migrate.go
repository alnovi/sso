package scripts

import "embed"

//go:embed migrations/*
var MigrateSchema embed.FS
