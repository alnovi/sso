package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
)

func (s *TestSuite) TestHttpOAuthForm() {
	session := &entity.Session{
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		query   map[string]string
		headers map[string]string
		cookies []*http.Cookie
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success show form",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			expCode: http.StatusOK,
			expBody: `<div id="app"></div>`,
		}, {
			name: "Success authorize with session_id",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusFound,
		}, {
			name: "Success show form with invalid session_id",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId("invalid", false),
			},
			expCode: http.StatusOK,
			expBody: `<div id="app"></div>`,
		}, {
			name: "Form response_type invalid",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "invalid",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			expCode: http.StatusBadRequest,
			expBody: `<div id="app"></div>`,
			expErr:  "Не валидный response-type",
		}, {
			name: "Json response_type invalid",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "invalid",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "Не валидный response-type",
		}, {
			name: "Form client_id invalid",
			query: map[string]string{
				"client_id":     "invalid",
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			expCode: http.StatusBadRequest,
			expBody: `<div id="app"></div>`,
			expErr:  "Клиент не найден",
		}, {
			name: "Json client_id invalid",
			query: map[string]string{
				"client_id":     "invalid",
				"response_type": "code",
				"redirect_uri":  s.config().CAdmin.Callback,
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "Клиент не найден",
		}, {
			name: "Form redirect_uri invalid",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  "invalid",
			},
			expCode: http.StatusBadRequest,
			expBody: `<div id="app"></div>`,
			expErr:  "Не валидный redirect-uri",
		}, {
			name: "Json redirect_uri invalid",
			query: map[string]string{
				"client_id":     s.config().CAdmin.Id,
				"response_type": "code",
				"redirect_uri":  "invalid",
			},
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			expCode: http.StatusBadRequest,
			expErr:  "Не валидный redirect-uri",
		},
	}

	ctrl := oauth.NewAuthController(s.app.Provider.OAuth(), s.app.Provider.Cookie())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)

			req := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
			s.applyHeaders(req, tc.headers)
			s.applyCookies(req, tc.cookies)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Form, c); err != nil {
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
