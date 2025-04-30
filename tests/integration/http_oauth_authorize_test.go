package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpOAuthAuthorize() {
	testCases := []struct {
		name      string
		query     map[string]string
		headers   map[string]string
		data      map[string]any
		expCode   int
		expBody   string
		expHeader map[string]string
		expErr    string
	}{
		{
			name: "Success authorize form",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationForm,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusFound,
			expHeader: map[string]string{
				"Location": s.config().CAdmin.Callback,
			},
		}, {
			name: "Success authorize json",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": s.config().UAdmin.Password,
				"remember": true,
			},
			expCode: http.StatusFound,
			expHeader: map[string]string{
				"Location": s.config().CAdmin.Callback,
			},
		}, {
			name: "Success authorize ajax",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type":     echo.MIMEApplicationJSON,
				"X-Requested-With": "XMLHttpRequest",
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusOK,
			expBody: s.config().CAdmin.Callback,
		}, {
			name: "Empty login authorize json",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    "",
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "login обязательное поле",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Invalid login authorize json",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    "invalid",
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "login должен быть email адресом",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Undefined login authorize json",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    "user@example.com",
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "пользователь не найден",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Empty password authorize json",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "password обязательное поле",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Invalid password authorize json",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": "qwerty",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "пароль не верный",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Invalid response_type authorize",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "invalid",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusBadRequest,
			expErr:  "Не валидный response-type",
		}, {
			name: "Invalid client_id authorize",
			query: map[string]string{
				"client_id":     "invalid",
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusBadRequest,
			expErr:  "Клиент не найден",
		}, {
			name: "Invalid redirect_uri authorize",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  "invalid",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"login":    s.config().UAdmin.Email,
				"password": s.config().UAdmin.Password,
			},
			expCode: http.StatusBadRequest,
			expErr:  "Не валидный redirect-uri",
		},
	}

	ms := []echo.MiddlewareFunc{
		middleware.TrailingSlash(),
	}
	ctrl := oauth.NewAuthController(s.app.Provider.OAuth(), s.app.Provider.Cookie())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)
			data := s.buildData(tc.headers["Content-Type"], tc.data)

			req := httptest.NewRequest(http.MethodPost, "/?"+query, strings.NewReader(data))
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err := s.sendToServer(ctrl.Authorize, c, ms...); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			for k, v := range tc.expHeader {
				s.Assert().Contains(rec.Header().Get(k), v, MsgNotAssertHeader)
			}

			s.Assert().Contains(rec.Body.String(), tc.expBody, MsgNotAssertBody)

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
