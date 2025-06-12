package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/alnovi/gomon/utils"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpOAuthResetPassword() {
	query := s.buildQuery(map[string]string{
		"client_id":    s.config().CAdmin.Id,
		"redirect_uri": s.config().CAdmin.Callback,
	})

	token, err := s.app.Provider.Token().ForgotPasswordToken(
		context.Background(),
		s.config().CAdmin.Id,
		s.config().UAdmin.Id,
		query,
		TestIP,
		TestAgent,
	)

	s.Require().NoError(err)

	testCases := []struct {
		name      string
		headers   map[string]string
		data      map[string]any
		expCode   int
		expBody   string
		expHeader map[string]string
		expErr    string
	}{
		{
			name: "Success",
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"token":    token.Hash,
				"password": "new-secret",
			},
			expCode: http.StatusFound,
			expHeader: map[string]string{
				"Location": "oauth/authorize?" + query,
			},
		}, {
			name: "Token is used",
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"token":    token.Hash,
				"password": "new-secret",
			},
			expCode: http.StatusBadRequest,
			expBody: "Токен не найден",
			expErr:  "token not found",
		}, {
			name: "Password is empty",
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"token":    token.Hash,
				"password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "password обязательное поле",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Password is small",
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"token":    token.Hash,
				"password": "123",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "password должен содержать минимум 5 символов",
			expErr:  "Unprocessable Entity",
		}, {
			name: "Password is large",
			headers: map[string]string{
				"Content-Type": echo.MIMEApplicationJSON,
			},
			data: map[string]any{
				"token":    token.Hash,
				"password": "123456789012345678901234567890",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "password должен содержать максимум 24 символа",
			expErr:  "Unprocessable Entity",
		},
	}

	ctrl := oauth.NewPasswordController(s.app.Provider.OAuth())

	for _, tc := range testCases {
		var user *entity.User

		s.Run(tc.name, func() {
			data := s.buildData(tc.headers["Content-Type"], tc.data)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.ResetPassword, c, middleware.TrailingSlash()); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			} else {
				user, err = s.app.Provider.StorageUsers().GetById(context.Background(), s.config().UAdmin.Id)
				s.Assert().NoError(err, MsgNotAssertError)
				s.Assert().Truef(utils.CompareHashPassword("new-secret", user.Password), "user password is failed")
			}

			for k, v := range tc.expHeader {
				s.Assert().Contains(rec.Header().Get(k), v, MsgNotAssertHeader)
			}

			s.Assert().Contains(rec.Body.String(), tc.expBody, MsgNotAssertBody)

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
