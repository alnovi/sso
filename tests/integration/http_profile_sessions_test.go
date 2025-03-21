package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestProfileSessions() {
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
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name:    "Success",
			expCode: http.StatusOK,
			expBody: []string{
				session.Id,
				session.Agent,
				session.Ip,
			},
		},
	}

	mdw := middleware.AuthBySession(s.app.Provider.Profile())
	ctrl := controller.NewProfileController(s.app.Provider.Profile(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
			req.Header.Set("User-Agent", TestAgent)
			req.AddCookie(s.app.Provider.Cookie().SessionId(session.Id, false))
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Sessions, c, mdw); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			for _, expBody := range tc.expBody {
				s.Assert().Contains(rec.Body.String(), expBody, MsgNotAssertBody)
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
