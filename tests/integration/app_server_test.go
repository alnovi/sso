package integration

import (
	"context"
	"time"

	"github.com/alnovi/sso/internal/app/server"
)

func (s *TestSuite) TestAppServerStart() {
	s.Require().NotPanics(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.NewApp(s.app.Provider.Config()).Start(ctx)
	})
}
