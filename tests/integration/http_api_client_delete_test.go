package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiClientDelete() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	_, access2, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleManager)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		client  string
		headers map[string]string
		expCode int
		expBody string
		expErr  string
	}{
		{
			name:   "Success",
			client: TestClient.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusOK,
			expBody: fmt.Sprintf(`"deleted_at":"%s`, time.Now().Format(time.DateOnly)),
		},
		{
			name:   "Success force",
			client: TestClient.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusOK,
			expBody: fmt.Sprintf(`"deleted_at":"%s`, time.Now().Format(time.DateOnly)),
		},
		{
			name:   "Success system client",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusNotFound,
			expErr:  "no results",
		},
		{
			name:   "Not found",
			client: "invalid",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusNotFound,
			expErr:  "no results",
		},
		{
			name:   "Unauthorized",
			client: TestClient.Id,
			headers: map[string]string{
				"User-Agent":   TestAgent,
				"Content-Type": "application/json",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized",
		},
		{
			name:   "Forbidden",
			client: TestClient.Id,
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
	ctrl := api.NewClientController(s.app.Provider.StorageClients())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/clients/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.client)

			if err = s.sendToServer(ctrl.Delete, c, mdws...); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			s.Assert().Contains(rec.Body.String(), tc.expBody, MsgNotAssertBody)

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
