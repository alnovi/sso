package integration

import (
	"context"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/crontask"
)

func (s *TestSuite) TestCronTaskDeleteSessionEmpty() {
	users := []string{s.config().UAdmin.Id, TestUser.Id}

	for _, user := range users {
		session := &entity.Session{Id: uuid.NewString(), UserId: user, Ip: TestIP, Agent: TestAgent}
		err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
		s.Require().NoError(err)
	}

	s.Run("delete session empty", func() {
		err := crontask.NewTaskDeleteSessionEmpty(s.app.Provider.Repository()).Handle()
		s.Require().NoError(err, "failed to delete empty session")

		sessionCount, err := s.app.Provider.Repository().SessionsCount(context.Background())
		s.Require().NoError(err)

		s.Require().Equal(0, sessionCount, "session count not equal")
	})
}
