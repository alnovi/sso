package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/transaport/http/controller/oauth"
)

func (s *TestSuite) TestOauthAuthorize() {
	testCases := []struct {
		name    string
		query   map[string]string
		data    map[string]any
		headers map[string]string
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success authorize json",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": s.app.Provider.Config().Admin.Password,
				"remember": true,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusFound,
			expBody: "",
			expErr:  "",
		},
		{
			name: "Success authorize ajax",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": s.app.Provider.Config().Admin.Password,
				"remember": true,
			},
			headers: map[string]string{
				"Content-Type":     echo.MIMEApplicationJSON,
				"X-Requested-With": "XMLHttpRequest",
			},
			expCode: http.StatusOK,
			expBody: s.app.Provider.Config().Client.Host,
			expErr:  "",
		},
		{
			name: "Success authorize form",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": s.app.Provider.Config().Admin.Password,
				"remember": true,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationForm,
			},
			expCode: http.StatusFound,
			expBody: "",
			expErr:  "",
		},
		{
			name: "Response type invalid",
			query: map[string]string{
				"response_type": "invalid",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": s.app.Provider.Config().Admin.Password,
				"remember": true,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  "response type invalid",
		},
		{
			name: "Client id invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     "invalid",
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": s.app.Provider.Config().Admin.Password,
				"remember": true,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  "client not found",
		},
		{
			name: "Redirect uri invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  "invalid",
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": s.app.Provider.Config().Admin.Password,
				"remember": true,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  "redirect url invalid",
		},
		{
			name: "User login empty",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    "",
				"password": s.app.Provider.Config().Admin.Password,
				"remember": false,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "User password empty",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": "",
				"remember": false,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "User login invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    "invalid",
				"password": s.app.Provider.Config().Admin.Password,
				"remember": false,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "User password invalid",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.app.Provider.Config().Client.Id,
				"redirect_uri":  s.app.Provider.Config().Client.Host,
			},
			data: map[string]any{
				"login":    s.app.Provider.Config().Admin.Email,
				"password": "invalid",
				"remember": false,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "",
			expErr:  "Unprocessable Entity",
		},
	}

	controller := oauth.NewAuthorizeController(s.app.Provider.OAuth(), s.app.Provider.Cookie())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)
			data := s.buildData(tc.headers["Content-Type"], tc.data)

			req := httptest.NewRequest(http.MethodPost, "/?"+query, strings.NewReader(data))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err := s.sendToServer(controller.Authorize, c); err != nil {
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
