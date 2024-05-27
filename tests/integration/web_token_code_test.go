package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

func (s *TestSuite) TestWebTokenCode() {
	ctx := context.Background()

	client, err := s.App.Provider.Repository().ClientByID(ctx, s.App.Provider.Config().Client.ProfileID)
	s.Require().NoError(err, "can't get client")

	activeCode := &entity.Token{
		Class:      entity.TokenClassCode,
		Hash:       "0000000001",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   client.ID,
		NotBefore:  time.Now().Add(-time.Second),
		Expiration: time.Now().Add(time.Minute),
	}

	expireCode := &entity.Token{
		Class:      entity.TokenClassCode,
		Hash:       "0000000002",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   client.ID,
		NotBefore:  time.Now().Add(-time.Minute),
		Expiration: time.Now().Add(-time.Second),
	}

	err = s.App.Provider.Repository().CreateToken(ctx, activeCode)
	s.Require().NoError(err, "can't create active code token")

	err = s.App.Provider.Repository().CreateToken(ctx, expireCode)
	s.Require().NoError(err, "can't create expire code token")

	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expBody string
		expErr  error
	}{
		{
			name: "Success code",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"code":          activeCode.Hash,
			},
			expCode: http.StatusOK,
			expBody: "",
			expErr:  nil,
		},
		{
			name: "Used code",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"code":          activeCode.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrTokenNotFound,
		},
		{
			name: "Invalid grant type",
			query: map[string]string{
				"grant_type":    "",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"code":          activeCode.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrUnsupportedGrantType,
		},
		{
			name: "Invalid client id",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     "",
				"client_secret": client.Secret,
				"code":          activeCode.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrClientNotFound,
		},
		{
			name: "Invalid client secret",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     client.ID,
				"client_secret": "",
				"code":          activeCode.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrClientNotFound,
		},
		{
			name: "Invalid code",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"code":          "",
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrTokenNotFound,
		},
		{
			name: "Expiration code",
			query: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"code":          expireCode.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrTokenExpired,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			q := make(url.Values)
			for k, v := range tc.query {
				q.Set(k, v)
			}

			req := httptest.NewRequest(http.MethodPost, "/oauth/token/?"+q.Encode(), nil)
			rec := httptest.NewRecorder()

			c := s.App.Server.NewContext(req, rec)

			if err = s.SendToServer(s.App.Provider.WebToken().Token, c); err != nil {
				s.Assert().ErrorIs(err, tc.expErr, "not assert error") //nolint:gosec
			}

			if tc.expBody != "" {
				s.Assert().Contains(rec.Body.String(), tc.expBody, "not assert body")
			}

			s.Assert().Equal(tc.expCode, rec.Code, "not assert code")
		})
	}
}
