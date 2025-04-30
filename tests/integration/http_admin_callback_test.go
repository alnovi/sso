package integration

import (
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller"
)

func (s *TestSuite) TestHttpAdminCallback() {
	session := &entity.Session{
		Id:     uuid.NewString(),
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
	s.Require().NoError(err)

	code, err := s.app.Provider.Token().CodeToken(context.Background(), session.Id, s.config().CAdmin.Id, s.config().UAdmin.Id)
	s.Require().NoError(err)

	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expBody string
		expErr  string
	}{
		{
			name: "Success",
			query: map[string]string{
				"code": code.Hash,
			},
			expCode: http.StatusFound,
		}, {
			name: "Invalid code",
			query: map[string]string{
				"code": "invalid",
			},
			expCode: http.StatusInternalServerError,
			expErr:  "token not found",
		},
	}

	ctrl := controller.NewAdminController(s.app.Provider.Admin(), s.app.Provider.Cookie(), nil)

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)

			req := httptest.NewRequest(http.MethodGet, "/?"+query, nil)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err := s.sendToServer(ctrl.Callback, c); err != nil {
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
