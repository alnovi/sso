package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpProfileLogout() {
	session := &entity.Session{
		Id:     uuid.NewString(),
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		session string
		expCode int
		expErr  string
	}{
		{
			name:    "Success",
			session: session.Id,
			expCode: http.StatusOK,
			expErr:  "",
		},
		{
			name:    "Unauthorized",
			session: uuid.NewString(),
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized: session not found",
		},
	}

	mdw := middleware.AuthBySession(s.app.Provider.Profile())
	ctrl := controller.NewProfileController(s.app.Provider.Profile(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("User-Agent", TestAgent)
			req.AddCookie(s.app.Provider.Cookie().SessionId(tc.session, false))
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Logout, c, mdw); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)

			_, err = s.app.Provider.Repository().SessionById(context.Background(), tc.session)

			s.Assert().ErrorIs(err, repository.ErrNoResult, "session not deleted")
		})
	}
}
