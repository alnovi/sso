package integration

import (
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/exception"
	"github.com/labstack/echo/v4"
)

func (s *TestSuite) TestWebAuthSignIn() {
	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expErr  error
	}{
		{
			name: "Success",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
			},
			expCode: http.StatusOK,
			expErr:  nil,
		}, {
			name: "Success with redirect uri",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
			},
			expCode: http.StatusOK,
			expErr:  nil,
		}, {
			name: "Empty or invalid client id",
			query: map[string]string{
				"response_type": "code",
				"client_id":     "",
			},
			expCode: http.StatusBadRequest,
			expErr:  echo.NewHTTPError(http.StatusBadRequest).SetInternal(exception.ErrClientNotFound),
		}, {
			name: "Empty or invalid response type",
			query: map[string]string{
				"response_type": "",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
			},
			expCode: http.StatusBadRequest,
			expErr:  echo.NewHTTPError(http.StatusBadRequest).SetInternal(exception.ErrUnsupportedGrantType),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/?"+s.BuildQuery(tc.query), nil)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			err := s.SendToServer(s.App.Provider.WebAuth().SignIn, c)

			s.Assert().Equal(tc.expErr, err, "not assert error")
			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}
