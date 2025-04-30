package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiSessionDelete() {
	session, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	session2, access2, _, err := s.accessTokens(s.config().CAdmin.Id, TestUser.Id, entity.RoleManager)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		session string
		headers map[string]string
		expCode int
		expErr  string
	}{
		{
			name:    "Success",
			session: session2.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
		},
		{
			name:    "Block current session",
			session: session.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusBadRequest,
			expErr:  "вы не можете удалить текущую сессию",
		},
		{
			name:    "Not found",
			session: "invalid",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusNotFound,
			expErr:  "no result",
		},
		{
			name: "Unauthorized",
			headers: map[string]string{
				"User-Agent":   TestAgent,
				"Content-Type": "application/json",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized",
		},
		{
			name: "Forbidden",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access2.Hash,
			},
			expCode: http.StatusForbidden,
			expErr:  "Forbidden",
		},
	}

	ms := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewSessionController(s.app.Provider.StorageSessions())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/sessions/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.session)

			if err = s.sendToServer(ctrl.Delete, c, ms...); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			if err == nil {
				_, err = s.app.Provider.Repository().SessionById(context.Background(), tc.session)
				s.Assert().EqualError(err, repository.ErrNoResult.Error(), MsgNotAssertError)
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
