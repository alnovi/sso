package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (s *TestSuite) TestWebAuthResetPasswordPage() {
	token := &entity.Token{
		Class:      entity.TokenClassResetPassword,
		Hash:       "0000000001",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   s.App.Provider.Config().Client.ProfileID,
		NotBefore:  time.Now().Add(-time.Second),
		Expiration: time.Now().Add(time.Minute),
	}

	err := s.App.Provider.Repository().CreateToken(context.Background(), token)
	s.Require().NoError(err, "can't create reset password token")

	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expErr  error
	}{
		{
			name:    "Success",
			query:   map[string]string{"hash": token.Hash},
			expCode: http.StatusOK,
			expErr:  nil,
		},
		{
			name:    "Empty hash",
			query:   map[string]string{"hash": ""},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrTokenNotFound,
		},
		{
			name:    "Invalid hash",
			query:   map[string]string{"hash": "invalid"},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrTokenNotFound,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/?"+s.BuildQuery(tc.query), nil)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			if err = s.SendToServer(s.App.Provider.WebAuth().ResetPasswordPage, c); err != nil {
				s.Assert().NotNil(tc.expErr, fmt.Sprintf("err is not nil: %s", err))
				s.Assert().ErrorAs(err, &tc.expErr, "not assert error") //nolint:gosec
			}

			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}

func (s *TestSuite) TestWebAuthResetPassword() {
	var err error
	var user *entity.User

	token1 := &entity.Token{
		Class:      entity.TokenClassResetPassword,
		Hash:       "0000000001",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   s.App.Provider.Config().Client.ProfileID,
		NotBefore:  time.Now().Add(-time.Second),
		Expiration: time.Now().Add(time.Minute),
	}

	token2 := &entity.Token{
		Class:      entity.TokenClassResetPassword,
		Hash:       "0000000002",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   s.App.Provider.Config().Client.ProfileID,
		NotBefore:  time.Now().Add(-time.Second),
		Expiration: time.Now().Add(time.Minute),
	}

	err = s.App.Provider.Repository().CreateToken(context.Background(), token1)
	s.Require().NoError(err, "can't create reset password token")

	err = s.App.Provider.Repository().CreateToken(context.Background(), token2)
	s.Require().NoError(err, "can't create reset password token")

	testCases := []struct {
		name    string
		mime    string
		form    map[string]string
		expPass string
		expCode int
		expErr  error
	}{
		{
			name: "Success form",
			mime: echo.MIMEApplicationForm,
			form: map[string]string{
				"hash":     token1.Hash,
				"password": "qwerty1",
			},
			expPass: "qwerty1",
			expCode: http.StatusFound,
			expErr:  nil,
		},
		{
			name: "Success json",
			mime: echo.MIMEApplicationJSON,
			form: map[string]string{
				"hash":     token2.Hash,
				"password": "qwerty2",
			},
			expPass: "qwerty2",
			expCode: http.StatusOK,
			expErr:  nil,
		},
		{
			name: "Token is used",
			mime: echo.MIMEApplicationJSON,
			form: map[string]string{
				"hash":     token1.Hash,
				"password": "qwerty",
			},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrTokenNotFound,
		},
		{
			name: "Empty hash",
			mime: echo.MIMEApplicationJSON,
			form: map[string]string{
				"hash":     "",
				"password": "qwerty",
			},
			expCode: http.StatusBadRequest,
			expErr:  exception.ErrTokenNotFound,
		},
		{
			name: "Empty password",
			mime: echo.MIMEApplicationJSON,
			form: map[string]string{
				"hash":     "invalid",
				"password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  &validator.ValidateError{},
		},
		{
			name: "Invalid password",
			mime: echo.MIMEApplicationJSON,
			form: map[string]string{
				"hash":     "invalid",
				"password": "qw",
			},
			expCode: http.StatusUnprocessableEntity,
			expErr:  &validator.ValidateError{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			data := s.BuildFormData(tc.mime, tc.form)

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
			req.Header.Set(echo.HeaderContentType, tc.mime)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			if err = s.SendToServer(s.App.WebAuth().ResetPassword, c); err != nil {
				s.Assert().NotNil(tc.expErr, fmt.Sprintf("err is not nil: %s", err))
				s.Assert().ErrorAs(err, &tc.expErr, fmt.Sprintf("not assert error: %s", err)) //nolint:gosec
			}

			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")

			if tc.expPass != "" {
				user, err = s.App.Repository().UserByID(context.Background(), s.App.Config().User.AdminID)
				s.Require().NoError(err, "fail get user")
				s.Assert().NoError(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tc.expPass)), "password incorrect")
			}
		})
	}
}
