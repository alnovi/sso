package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/cookie"
	"github.com/alnovi/sso/internal/transport/http/controller"
	"github.com/alnovi/sso/internal/transport/http/middleware"
)

func (s *TestSuite) TestHttpProfileHome() {
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
		query   map[string]string
		headers map[string]string
		cookies []*http.Cookie
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success with query",
			headers: map[string]string{
				"User-Agent": TestAgent,
			},
			query: map[string]string{
				cookie.SessionId: session.Id,
			},
			expCode: http.StatusFound,
		}, {
			name: "Success with cookie",
			headers: map[string]string{
				"User-Agent": TestAgent,
			},
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(session.Id, false),
			},
			expCode: http.StatusOK,
			expBody: "SSO | Профиль пользователя",
		}, {
			name:    "Unauthorized",
			expCode: http.StatusUnauthorized,
			expBody: "SSO | Ошибка",
			expErr:  "session not found",
		}, {
			name: "Unauthorized with query",
			query: map[string]string{
				cookie.SessionId: uuid.NewString(),
			},
			expCode: http.StatusUnauthorized,
			expBody: "SSO | Ошибка",
			expErr:  "session not found",
		}, {
			name: "Unauthorized with cookie",
			cookies: []*http.Cookie{
				s.app.Provider.Cookie().SessionId(uuid.NewString(), false),
			},
			expCode: http.StatusUnauthorized,
			expBody: "SSO | Ошибка",
			expErr:  "session not found",
		},
	}

	mdw := middleware.AuthBySession(s.app.Provider.Profile())
	ctrl := controller.NewProfileController(s.app.Provider.Profile(), s.app.Provider.Cookie(), mdw)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)

			req := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
			s.applyHeaders(req, tc.headers)
			s.applyCookies(req, tc.cookies)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Home, c, mdw); err != nil {
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
