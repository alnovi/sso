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

func (s *TestSuite) TestHttpProfileUpdateUser() {
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
		expBody []string
		expErr  string
	}{
		{
			name: "Success",
			data: map[string]any{
				"name":  "Иванов Иван Иванович",
				"email": "ivan@example.com",
			},
			expCode: http.StatusOK,
			expBody: []string{
				s.config().UAdmin.Id,
				"Иванов Иван Иванович",
				"ivan@example.com",
			},
		}, {
			name: "Invalid validation",
			data: map[string]any{
				"name":  "",
				"email": "ivan@example",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				"name обязательное поле",
				"email должен быть email адресом",
			},
			expErr: "Unprocessable Entity",
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

			if err = s.sendToServer(ctrl.UpdateUser, c, mdw); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			for _, expBody := range tc.expBody {
				s.Assert().Contains(rec.Body.String(), expBody, MsgNotAssertBody)
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
