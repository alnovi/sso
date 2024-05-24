package integration

import (
	"net/http"
	"net/http/httptest"
	"net/url"

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
				"redirect_uri":  "",
			},
			expCode: http.StatusOK,
			expErr:  nil,
		}, {
			name: "Success with redirect uri",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/profile/callback",
			},
			expCode: http.StatusOK,
			expErr:  nil,
		}, {
			name: "Empty or invalid client id",
			query: map[string]string{
				"response_type": "code",
				"client_id":     "",
				"redirect_uri":  "",
			},
			expCode: http.StatusBadRequest,
			expErr:  echo.NewHTTPError(http.StatusBadRequest).SetInternal(exception.ErrClientNotFound),
		}, {
			name: "Empty or invalid response type",
			query: map[string]string{
				"response_type": "",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "",
			},
			expCode: http.StatusBadRequest,
			expErr:  echo.NewHTTPError(http.StatusBadRequest).SetInternal(exception.ErrUnsupportedGrantType),
		}, {
			name: "Invalid redirect uri",
			query: map[string]string{
				"response_type": "code",
				"client_id":     s.App.Provider.Config().Client.ProfileID,
				"redirect_uri":  "/callback/invalid",
			},
			expCode: http.StatusBadRequest,
			expErr:  echo.NewHTTPError(http.StatusBadRequest).SetInternal(exception.ErrClientNotFound),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			q := make(url.Values)
			for k, v := range tc.query {
				q.Set(k, v)
			}

			req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			err := s.SendToServer(s.App.Provider.WebAuth().SignIn, c)

			s.Assert().Equal(tc.expErr, err, "not assert error")
			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}
