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

func (s *TestSuite) TestHttpApiUserClients() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	_, access2, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleManager)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		user    string
		headers map[string]string
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name: "Success admin",
			user: s.config().UAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusOK,
			expBody: []string{
				s.config().CAdmin.Id,
				s.config().CAdmin.Name,
				entity.RoleAdmin,
			},
		}, {
			name: "Success manager",
			user: TestUser.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusOK,
			expBody: []string{
				TestClient.Id,
				TestClient.Name,
				entity.RoleManager,
			},
		}, {
			name: "Unauthorized",
			user: s.config().UAdmin.Id,
			headers: map[string]string{
				"User-Agent":   TestAgent,
				"Content-Type": "application/json",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized",
		}, {
			name: "Forbidden",
			user: s.config().UAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access2.Hash),
			},
			expCode: http.StatusForbidden,
			expErr:  "Forbidden",
		},
	}

	mdws := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewUserController(s.app.Provider.StorageUsers(), s.app.Provider.StorageRoles())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/users/:id/clients")
			c.SetParamNames("id")
			c.SetParamValues(tc.user)

			if err = s.sendToServer(ctrl.Clients, c, mdws...); err != nil {
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
