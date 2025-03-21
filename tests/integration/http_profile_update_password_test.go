package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpProfileUpdatePassword() {
	session := &entity.Session{
		Id:     uuid.NewString(),
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		data    map[string]any
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success",
			data: map[string]any{
				"old_password": s.config().UAdmin.Password,
				"new_password": "new_password",
			},
			expCode: http.StatusOK,
		},
		{
			name: "Incorrect old password",
			data: map[string]any{
				"old_password": s.config().UAdmin.Password,
				"new_password": "new_password",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "Пароль не верный",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Empty old password",
			data: map[string]any{
				"old_password": "",
				"new_password": "new_password",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "old_password обязательное поле",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Short old password",
			data: map[string]any{
				"old_password": "1234",
				"new_password": "new_password",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "old_password должен содержать минимум 5 символов",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Empty new password",
			data: map[string]any{
				"old_password": s.config().UAdmin.Password,
				"new_password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "new_password обязательное поле",
			expErr:  "Unprocessable Entity",
		},
		{
			name: "Short new password",
			data: map[string]any{
				"old_password": s.config().UAdmin.Password,
				"new_password": "1234",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: "new_password должен содержать минимум 5 символов",
			expErr:  "Unprocessable Entity",
		},
	}

	mdw := middleware.AuthBySession(s.app.Provider.Profile())
	ctrl := controller.NewProfileController(s.app.Provider.Profile(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			data := s.buildDataJson(tc.data)

			req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(data))
			req.Header.Set("Content-Type", echo.MIMEApplicationJSON)
			req.Header.Set("User-Agent", TestAgent)
			req.AddCookie(s.app.Provider.Cookie().SessionId(session.Id, false))
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.UpdatePassword, c, mdw); err != nil {
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
