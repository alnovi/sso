package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (s *TestSuite) TestWebAuthAuthorize() {
	testCases := []struct {
		name    string
		mime    string
		query   map[string]string
		form    map[string]string
		expCode int
		expBody string
		expErr  error
	}{
		{
			name: "Success form",
			mime: echo.MIMEApplicationForm,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "",
			},
			form: map[string]string{
				"login":    "admin@example.ru",
				"password": "admin",
			},
			expCode: http.StatusFound,
			expErr:  nil,
		},
		{
			name: "Success form with redirect uri",
			mime: echo.MIMEApplicationForm,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "admin@example.ru",
				"password": "admin",
			},
			expCode: http.StatusFound,
			expErr:  nil,
		},
		{
			name: "Success json",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "",
			},
			form: map[string]string{
				"login":    "admin@example.ru",
				"password": "admin",
			},
			expCode: http.StatusOK,
			expBody: `"location":"/profile/callback?code=`,
			expErr:  nil,
		},
		{
			name: "Success json with redirect uri",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "admin@example.ru",
				"password": "admin",
			},
			expCode: http.StatusOK,
			expBody: `"location":"/profile/callback?code=`,
			expErr:  nil,
		},
		{
			name: "Empty fields form",
			mime: echo.MIMEApplicationForm,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "",
				"password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  &validator.ValidateError{},
		},
		{
			name: "Empty fields json",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "",
				"password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"validate":{"login":"login обязательное поле","password":"password обязательное поле"}`,
			expErr:  &validator.ValidateError{},
		},
		{
			name: "Login not email",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "admin",
				"password": "admin",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"validate":{"login":"login должен быть email адресом"}`,
			expErr:  &validator.ValidateError{},
		},
		{
			name: "User not found",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "example@example.ru",
				"password": "admin",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"validate":{"login":"Пользователь не найден"}`,
			expErr:  &validator.ValidateError{},
		},
		{
			name: "Password invalid",
			mime: echo.MIMEApplicationJSON,
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			form: map[string]string{
				"login":    "admin@example.ru",
				"password": "secret",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: `"validate":{"password":"Не верный пароль"}`,
			expErr:  &validator.ValidateError{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			var data string

			query := make(url.Values)
			for k, v := range tc.query {
				query.Set(k, v)
			}

			if tc.mime == echo.MIMEApplicationForm {
				form := make(url.Values)
				for k, v := range tc.form {
					form.Set(k, v)
				}
				data = form.Encode()
			} else {
				form, _ := json.Marshal(tc.form)
				data = string(form)
			}

			req := httptest.NewRequest(http.MethodPost, "/?"+query.Encode(), strings.NewReader(data))
			req.Header.Set(echo.HeaderContentType, tc.mime)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			if err := s.SendToServer(s.App.Provider.WebAuth().Authorize, c); err != nil {
				s.Assert().ErrorAs(err, &tc.expErr, "not assert error") //nolint:gosec
			}

			if tc.expBody != "" {
				s.Assert().Contains(rec.Body.String(), tc.expBody, "not assert body")
			}

			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}
