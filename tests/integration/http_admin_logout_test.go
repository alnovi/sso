package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpAdminLogout() {
	session := &entity.Session{
		Id:     uuid.NewString(),
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
	s.Require().NoError(err)

	access, err := s.app.Provider.Token().AccessToken(context.Background(), session.Id, s.config().CAdmin.Id, s.config().UAdmin.Id, entity.RoleAdmin)
	s.Require().NoError(err)

	refresh, err := s.app.Provider.Token().RefreshToken(context.Background(), session.Id, s.config().CAdmin.Id, s.config().UAdmin.Id, time.Now())
	s.Require().NoError(err)

	testCase := struct {
		name    string
		cookies []*http.Cookie
		expCode int
	}{
		name: "Success",
		cookies: []*http.Cookie{
			s.app.Provider.Cookie().SessionId(session.Id, false),
			s.app.Provider.Cookie().AccessToken(access),
			s.app.Provider.Cookie().RefreshToken(refresh),
		},
		expCode: http.StatusOK,
	}

	mdw := middleware.Token(s.app.Provider.OAuth(), s.app.Provider.Cookie(), s.app.Provider.Config().CAdmin.Id, s.app.Provider.Config().CAdmin.Secret)
	ctrl := controller.NewAdminController(s.app.Provider.Admin(), s.app.Provider.Cookie(), mdw)

	s.Run("Success", func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		s.applyCookies(req, testCase.cookies)
		rec := httptest.NewRecorder()

		c := s.app.HttpServer.NewContext(req, rec)

		if err = s.sendToServer(ctrl.Logout, c, mdw); err != nil {
			s.Assert().NoError(err, MsgNotAssertError)
		}

		s.Assert().Equal(testCase.expCode, rec.Code, MsgNotAssertCode)

		removeCookies := rec.Header().Values("Set-Cookie")
		s.Assert().Len(removeCookies, 3, "count cookies not attempts")

		_, err = s.app.Provider.Token().ValidateRefreshToken(context.Background(), refresh.Hash)
		s.Assert().Error(err)
	})
}
