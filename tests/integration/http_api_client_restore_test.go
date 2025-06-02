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

func (s *TestSuite) TestHttpApiClientRestore() {
	_, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	_, err = s.app.Provider.StorageClients().Delete(context.Background(), TestClient.Id)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		client  string
		headers map[string]string
		expCode int
		expBody string
		expErr  string
	}{
		{
			name:   "Success",
			client: TestClient.Id,
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusOK,
			expBody: `"deleted_at":null`,
		}, {
			name:   "Not found",
			client: "invalid",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", access.Hash),
			},
			expCode: http.StatusNotFound,
			expErr:  "no results",
		},
	}

	mdws := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewClientController(s.app.Provider.StorageClients())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)
			c.SetPath("/api/clients/:id/restore")
			c.SetParamNames("id")
			c.SetParamValues(tc.client)

			if err = s.sendToServer(ctrl.Restore, c, mdws...); err != nil {
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
