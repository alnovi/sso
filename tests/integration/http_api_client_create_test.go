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

func (s *TestSuite) TestHttpApiClientCreate() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	_, access2, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleManager)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		data    map[string]any
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusOK,
			expBody: "integration-test-client",
		},
		{
			name: "Invalid id empty",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"id":"id обязательное поле"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid id is not lowercase",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "Test~Client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"id":"Значение может содержать только буквы (в нижнем регистре), цифры и дефис"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid id is use",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "test-client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"id":"Такое значение уже занято"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid id length min",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "id",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"id":"id должен содержать минимум 3 символа"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid id length max",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "client-id-is-big-length-1234567890",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"id":"id должен содержать максимум 30 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid data",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       123,
				"name":     123,
				"icon":     nil,
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusBadRequest,
			expErr:  "Bad Request",
		},
		{
			name: "Invalid name empty",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"name":"name обязательное поле"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid name length min",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "1234",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"name":"name должен содержать минимум 5 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid name length max",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "1234567890-1234567890-1234567890-1234567890-1234567890",
				"icon":     "https://example.com/icon.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"name":"name должен содержать максимум 50 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid icon is not url",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "Client auto create with testing",
				"icon":     "example.png",
				"callback": "https://example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"icon":"icon должен быть URI"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid callback is empty",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/example.png",
				"callback": "",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"callback":"callback обязательное поле"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid callback is not url",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/example.png",
				"callback": "example.com/callback",
				"secret":   "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"callback":"callback должен быть URL"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid secret is small",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/example.png",
				"callback": "https://example.com/callback",
				"secret":   "1234",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"secret":"secret должен содержать минимум 5 символов"`,
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Invalid secret is big",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			data: map[string]any{
				"id":       "integration-test-client",
				"name":     "Client auto create with testing",
				"icon":     "https://example.com/example.png",
				"callback": "https://example.com/callback",
				"secret":   "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"secret":"secret должен содержать максимум 100 символов"`,
			expErr:  "Unprocessable Entity",
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
			data := s.buildDataJson(tc.data)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Create, c, mdws...); err != nil {
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
