package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiUserCreate() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		data    map[string]any
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name: "Success",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  echo.MIMEApplicationJSON,
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "Иванов Иван Иванович",
				"email":    "ivan@example.com",
				"password": TestSecret,
			},
			expCode: http.StatusOK,
			expBody: []string{
				"Иванов Иван Иванович",
				"ivan@example.com",
			},
		},
		{
			name: "Invalid name is empty",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "",
				"email":    "ivan@example.com",
				"password": TestSecret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"name":"name обязательное поле"`,
			},
			expErr: "Unprocessable Entity",
		},
		{
			name: "Invalid email is empty",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "Иванов Иван Иванович",
				"email":    "",
				"password": TestSecret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"email":"email обязательное поле"`,
			},
			expErr: "Unprocessable Entity",
		},
		{
			name: "Invalid email is use",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "Иванов Иван Иванович",
				"email":    "ivan@example.com",
				"password": TestSecret,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"email":"Такое значение уже занято"`,
			},
			expErr: "Unprocessable Entity",
		},
		{
			name: "Invalid password is empty",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "Иванов Иван Иванович",
				"email":    "ivan@example.com",
				"password": "",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"password":"password обязательное поле"`,
			},
			expErr: "Unprocessable Entity",
		},
		{
			name: "Invalid password min length",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "Иванов Иван Иванович",
				"email":    "ivan@example.com",
				"password": "12",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"password":"password должен содержать минимум 5 символов"`,
			},
			expErr: "Unprocessable Entity",
		},
		{
			name: "Invalid password max length",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			data: map[string]any{
				"name":     "Иванов Иван Иванович",
				"email":    "ivan@example.com",
				"password": "123456789012345678901234567890",
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"password":"password должен содержать максимум 24 символа"`,
			},
			expErr: "Unprocessable Entity",
		},
	}

	mdws := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewUserController(s.app.Provider.StorageUsers(), s.app.Provider.StorageRoles())

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

			if len(tc.expBody) > 0 {
				for _, body := range tc.expBody {
					s.Assert().Contains(rec.Body.String(), body, MsgNotAssertBody)
				}
			}

			s.Assert().Equal(tc.expCode, rec.Code, MsgNotAssertCode)
		})
	}
}
