package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

func (s *TestSuite) TestWebTokenRefresh() {
	ctx := context.Background()

	client, err := s.App.Provider.Repository().ClientByID(ctx, s.App.Provider.Config().Client.ProfileID)
	s.Require().NoError(err, "can't get client")

	activeToken := &entity.Token{
		Class:      entity.TokenClassRefresh,
		Hash:       "0000000001",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   client.ID,
		NotBefore:  time.Now().Add(-time.Minute),
		Expiration: time.Now().Add(time.Minute),
	}

	notBeforeToken := &entity.Token{
		Class:      entity.TokenClassRefresh,
		Hash:       "0000000002",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   client.ID,
		NotBefore:  time.Now().Add(time.Minute),
		Expiration: time.Now().Add(time.Hour),
	}

	expirationToken := &entity.Token{
		Class:      entity.TokenClassRefresh,
		Hash:       "0000000003",
		UserID:     s.App.Provider.Config().User.AdminID,
		ClientID:   client.ID,
		NotBefore:  time.Now().Add(-time.Hour),
		Expiration: time.Now().Add(-time.Minute),
	}

	err = s.App.Provider.Repository().CreateToken(ctx, activeToken)
	s.Require().NoError(err, "can't create active token refresh")

	err = s.App.Provider.Repository().CreateToken(ctx, notBeforeToken)
	s.Require().NoError(err, "can't create not before token refresh")

	err = s.App.Provider.Repository().CreateToken(ctx, expirationToken)
	s.Require().NoError(err, "can't create expiration token refresh")

	testCases := []struct {
		name    string
		query   map[string]string
		expCode int
		expBody string
		expErr  error
	}{
		{
			name: "Success refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"refresh_token": activeToken.Hash,
			},
			expCode: http.StatusOK,
			expBody: "",
			expErr:  nil,
		},
		{
			name: "Used refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"refresh_token": activeToken.Hash,
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
				"refresh_token": activeToken.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrUnsupportedGrantType,
		},
		{
			name: "Invalid client id",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     "",
				"client_secret": client.Secret,
				"refresh_token": activeToken.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrClientNotFound,
		},
		{
			name: "Invalid client secret",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     client.ID,
				"client_secret": "",
				"refresh_token": activeToken.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrClientNotFound,
		},
		{
			name: "Invalid refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"refresh_token": "",
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrTokenNotFound,
		},
		{
			name: "Not before refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"refresh_token": notBeforeToken.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrTokenUsedBeforeIssued,
		},
		{
			name: "Expiration refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     client.ID,
				"client_secret": client.Secret,
				"refresh_token": expirationToken.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "",
			expErr:  exception.ErrTokenExpired,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/oauth/token/?"+s.BuildQuery(tc.query), nil)
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
