package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiClientUpdate() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		client  string
		headers map[string]string
		data    map[string]any
		expCode int
		expBody string
		expErr  string
	}{
		{
			name:   "Success",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     s.config().CAdmin.Name,
				"icon":     nil,
				"callback": s.config().CAdmin.Callback,
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusOK,
			expBody: s.config().CAdmin.Id,
		},
		{
			name:   "Invalid name empty",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     "",
				"icon":     nil,
				"callback": s.config().CAdmin.Callback,
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"name":"name обязательное поле"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid name length min",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     "1234",
				"icon":     nil,
				"callback": s.config().CAdmin.Callback,
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"name":"name должен содержать минимум 5 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid name length max",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     "1234567890-1234567890-1234567890-1234567890-1234567890",
				"icon":     nil,
				"callback": s.config().CAdmin.Callback,
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"name":"name должен содержать максимум 50 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid icon is not url",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     s.config().CAdmin.Name,
				"icon":     "example.png",
				"callback": s.config().CAdmin.Callback,
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"icon":"icon должен быть URI"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid callback is empty",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     s.config().CAdmin.Name,
				"icon":     nil,
				"callback": "",
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"callback":"callback обязательное поле"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid callback is not url",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     s.config().CAdmin.Name,
				"icon":     nil,
				"callback": "example.com/callback",
				"secret":   s.config().CAdmin.Secret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"callback":"callback должен быть URI"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid secret is small",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     s.config().CAdmin.Name,
				"icon":     nil,
				"callback": s.config().CAdmin.Callback,
				"secret":   "1234",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"secret":"secret должен содержать минимум 5 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name:   "Invalid secret is big",
			client: s.config().CAdmin.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"name":     s.config().CAdmin.Name,
				"icon":     nil,
				"callback": s.config().CAdmin.Callback,
				"secret":   "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"secret":"secret должен содержать максимум 100 символов"`,
			expErr:  "Unprocessable Entity",
		},
	}

	ms := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewClientController(s.app.Provider.StorageClients())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			data := s.buildDataJson(tc.data)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/clients/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.client)

			if err = s.sendToServer(ctrl.Update, c, ms...); err != nil {
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
