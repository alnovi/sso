package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiSessionGet() {
	session, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	session2, _, _, err := s.accessTokens(s.config().CAdmin.Id, TestUser.Id, entity.RoleManager)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		session string
		headers map[string]string
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name:    "Success is current",
			session: session.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
			expBody: []string{
				`"is_current":true`,
				fmt.Sprintf(`"ip":"%s"`, session.Ip),
				fmt.Sprintf(`"agent":"%s"`, session.Agent),
				fmt.Sprintf(`"name":"%s"`, s.config().UAdmin.Name),
			},
		},
		{
			name:    "Success is not current",
			session: session2.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
			expBody: []string{
				`"is_current":false`,
				fmt.Sprintf(`"ip":"%s"`, session2.Ip),
				fmt.Sprintf(`"agent":"%s"`, session2.Agent),
				fmt.Sprintf(`"name":"%s"`, TestUser.Name),
			},
		},
	}

	ms := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewSessionController(s.app.Provider.StorageSessions())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/sessions/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.session)

			if err = s.sendToServer(ctrl.Get, c, ms...); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			if len(tc.expBody) > 0 {
				for _, body := range tc.expBody {
					s.Assert().Contains(rec.Body.String(), body, MsgNotAssertBody)
				}
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
