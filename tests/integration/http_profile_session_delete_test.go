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

func (s *TestSuite) TestHttpProfileSessionDelete() {
	sessionAdmin := &entity.Session{
		Id:     uuid.NewString(),
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	sessionTest := &entity.Session{
		Id:     uuid.NewString(),
		UserId: TestUser.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), sessionAdmin)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().SessionCreate(context.Background(), sessionTest)
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		sessionId string
		expCode   int
		expErr    string
	}{
		{
			name:      "Not user session",
			sessionId: sessionTest.Id,
			expCode:   http.StatusBadRequest,
			expErr:    "session not found",
		},
		{
			name:      "Success",
			sessionId: sessionAdmin.Id,
			expCode:   http.StatusOK,
		},
	}

	mdw := middleware.AuthBySession(s.app.Provider.Profile())
	ctrl := controller.NewProfileController(s.app.Provider.Profile(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
			req.Header.Set("User-Agent", TestAgent)
			req.AddCookie(s.app.Provider.Cookie().SessionId(sessionAdmin.Id, false))
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tc.sessionId)

			if err = s.sendToServer(ctrl.SessionDelete, c, mdw); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
