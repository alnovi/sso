package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/alnovi/gomon/utils"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiUserUpdateRole() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		user    string
		client  string
		role    *string
		headers map[string]string
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name:   "Success",
			user:   TestUser.Id,
			client: TestClient.Id,
			role:   utils.Point(entity.RoleGuest),
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
		},
		{
			name:   "Success remove role",
			user:   TestUser.Id,
			client: TestClient.Id,
			role:   nil,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
		},
		{
			name:   "Invalid role",
			user:   TestUser.Id,
			client: TestClient.Id,
			role:   utils.Point("invalid"),
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusUnprocessableEntity,
			expBody: []string{
				`"error":"Ошибка ввода данных"`,
				`"role":"role должен быть одним из [guest user manager admin]"`,
			},
			expErr: "Unprocessable Entity",
		},
	}

	ms := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewUserController(s.app.Provider.StorageUsers(), s.app.Provider.StorageRoles())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			data := s.buildDataJson(map[string]any{"role": tc.role})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(data))
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/users/:uid/clients/:cid")
			c.SetParamNames("uid", "cid")
			c.SetParamValues(tc.user, tc.client)

			if err = s.sendToServer(ctrl.UpdateRole, c, ms...); err != nil {
				if tc.expErr != "" {
					s.Assert().ErrorContains(err, tc.expErr, MsgNotAssertError)
				} else {
					s.Assert().NoError(err, MsgNotAssertError)
				}
			}

			if err == nil && tc.role != nil {
				role, _ := s.app.Provider.Repository().Role(context.Background(), tc.client, tc.user)
				s.Assert().Equal(role.Role, *tc.role, "not assert role")
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
