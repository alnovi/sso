package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/api"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpApiSessionList() {
	session, access, _, err := s.accessTokens(s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		headers map[string]string
		expCode int
		expBody []string
		expErr  string
	}{
		{
			name: "Success",
			headers: map[string]string{
				"User-Agent":    TestAgent,
				"Content-Type":  "application/json",
				"Authorization": access.Hash,
			},
			expCode: http.StatusOK,
			expBody: []string{
				fmt.Sprintf(`"ip":"%s"`, session.Ip),
				fmt.Sprintf(`"agent":"%s"`, session.Agent),
				fmt.Sprintf(`"name":"%s"`, s.config().UAdmin.Name),
			},
		},
	}

	ms := []echo.MiddlewareFunc{
		middleware.Auth(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret),
		middleware.RoleWeight(entity.RoleAdminWeight),
	}
	ctrl := api.NewSessionController(s.app.Provider.StorageSessions())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			s.applyHeaders(req, tc.headers)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.List, c, ms...); err != nil {
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
