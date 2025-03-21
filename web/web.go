package web

import "embed"

//go:embed public/* out/html/*.html out/assets/*
var StaticFS embed.FS
