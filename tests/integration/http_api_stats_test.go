package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/token"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiStats() {
	_, accessAdmin, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	_, accessExp, refreshExp, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin, token.WithAccessExpiresAt(time.Now().Add(-time.Minute)))
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		cookies []*http.Cookie
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success with access",
			headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", accessAdmin.Hash),
			},
			cookies: []*http.Cookie{},
			expCode: http.StatusOK,
			expBody: `{"users":2,"clients":2,"sessions":2}`,
		},
		{
			name: "Success with refresh",
			headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": accessExp.Hash,
				"Refresh-Token": refreshExp.Hash,
			},
			cookies: []*http.Cookie{},
			expCode: http.StatusOK,
			expBody: `{"users":2,"clients":2,"sessions":2}`,
		},
	}

	mdws := []echo.MiddlewareFunc{
		middleware.RequestLogger(s.app.Provider.Logger()),
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewStatsController(s.app.Provider.Stats())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			s.applyCookies(req, tc.cookies)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Stats, c, mdws...); err != nil {
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
