package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiUserRestore() {
	_, err := s.app.Provider.StorageUsers().Delete(context.Background(), TestUser.Id)
	s.Require().NoError(err)

	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	_, access2, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleManager)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		user    string
		headers map[string]string
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name: "Success",
			user: TestUser.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
			expBody: []string{
				fmt.Sprintf(`"id":"%s"`, TestUser.Id),
				fmt.Sprintf(`"email":"%s"`, TestUser.Email),
				`"deleted_at":null`,
			},
		},
		{
			name: "Not found",
			user: "invalid",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusNotFound,
			expErr:  "no results",
		},
		{
			name: "Unauthorized",
			user: TestUser.Id,
			headers: map[string]string{
				"User-Agent":   TestAgent,
				"Content-Type": "application/json",
			},
			expCode: http.StatusUnauthorized,
			expErr:  "Unauthorized",
		},
		{
			name: "Forbidden",
			user: TestUser.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access2.Hash,
			},
			expCode: http.StatusForbidden,
			expErr:  "Forbidden",
		},
	}

	mdws := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewUserController(s.app.Provider.StorageUsers(), s.app.Provider.StorageRoles())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/users/:id/restore")
			c.SetParamNames("id")
			c.SetParamValues(tc.user)

			if err = s.sendToServer(ctrl.Restore, c, mdws...); err != nil {
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
