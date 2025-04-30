package integration

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/crontask"
)

func (s *TestSuite) TestCronTaskDeleteTokenExpired() {
	hashes := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}

	for i, hash := range hashes {
		token := &entity.Token{
			Class:      entity.TokenClassRefresh,
			Hash:       hash,
			NotBefore:  time.Now(),
			Expiration: time.Now(),
		}
		err := s.app.Provider.Repository().TokenCreate(context.Background(), token)
		s.Require().NoErrorf(err, "fail create token #%d: %s", i, err)
	}

	s.Run("delete token expired", func() {
		err := crontask.NewTaskDeleteTokenExpired(s.app.Provider.Repository()).Handle()
		s.Require().NoError(err, "failed to delete token expired")

		for i, hash := range hashes {
			_, err = s.app.Provider.Repository().TokenByHash(context.Background(), hash)
			s.Assert().ErrorIsf(err, repository.ErrNoResult, "token #%d not deleted", i)
		}
	})
}
