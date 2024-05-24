package integration

import (
	"context"
	"time"

	"github.com/alnovi/sso/internal/app/server"
)

func (s *TestSuite) TestAppServer() {
	s.NotPanics(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := server.New(s.App.Provider.Config()).Start(ctx)
		s.NoError(err)
	})
}
