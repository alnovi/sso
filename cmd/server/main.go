package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/alnovi/sso/internal/app/server"
	"github.com/alnovi/sso/internal/config"
	"github.com/alnovi/sso/pkg/configure"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := &config.Config{}

	must(configure.ParseEnv(ctx, cfg))
	must(server.New(cfg).Start(ctx))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
