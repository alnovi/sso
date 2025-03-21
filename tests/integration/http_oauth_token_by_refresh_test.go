package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
	"github.com/alnovi/sso/pkg/rand"
)

func (s *TestSuite) TestHttpOAuthTokenByRefresh() {
	session := &entity.Session{
		Id:     uuid.NewString(),
		UserId: s.config().UAdmin.Id,
		Ip:     TestIP,
		Agent:  TestAgent,
	}

	code := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassCode,
		Hash:       rand.Base62(entity.TokenCodeCost),
		SessionId:  &session.Id,
		UserId:     &s.config().UAdmin.Id,
		ClientId:   &s.config().CAdmin.Id,
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(entity.TokenCodeTTL),
	}

	refresh := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(entity.TokenRefreshCost),
		SessionId:  &session.Id,
		UserId:     &s.config().UAdmin.Id,
		ClientId:   &s.config().CAdmin.Id,
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(entity.TokenRefreshTTL),
	}

	refreshNotSession := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(entity.TokenRefreshCost),
		SessionId:  nil,
		UserId:     &s.config().UAdmin.Id,
		ClientId:   &s.config().CAdmin.Id,
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(entity.TokenRefreshTTL),
	}

	refreshNotBefore := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(entity.TokenRefreshCost),
		SessionId:  &session.Id,
		UserId:     &s.config().UAdmin.Id,
		ClientId:   &s.config().CAdmin.Id,
		NotBefore:  code.Expiration,
		Expiration: time.Now().Add(entity.TokenRefreshTTL),
	}

	refreshExpiration := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(entity.TokenRefreshCost),
		SessionId:  &session.Id,
		UserId:     &s.config().UAdmin.Id,
		ClientId:   &s.config().CAdmin.Id,
		NotBefore:  time.Now().Add(-entity.TokenRefreshTTL),
		Expiration: time.Now(),
	}

	err := s.app.Provider.Repository().SessionCreate(context.Background(), session)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().TokenCreate(context.Background(), code)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().TokenCreate(context.Background(), refresh)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().TokenCreate(context.Background(), refreshNotSession)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().TokenCreate(context.Background(), refreshNotBefore)
	s.Require().NoError(err)

	err = s.app.Provider.Repository().TokenCreate(context.Background(), refreshExpiration)
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
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refresh.Hash,
			},
			expCode: http.StatusOK,
			expBody: "access_token",
		}, {
			name: "Reuse refresh token",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refresh.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "token not found",
			expErr:  "token not found",
		}, {
			name: "Invalid token not before refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refreshNotBefore.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "token not found",
			expErr:  "token not found",
		}, {
			name: "Invalid token expiration refresh",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refreshExpiration.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "token not found",
			expErr:  "token not found",
		}, {
			name: "Invalid token with empty session",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refreshNotSession.Hash,
			},
			expCode: http.StatusInternalServerError,
			expBody: "Ошибка сервера",
			expErr:  "session not found",
		}, {
			name: "Invalid token class",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refresh.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "token not found",
			expErr:  "token not found",
		}, {
			name: "Invalid grant_type",
			query: map[string]string{
				"grant_type":    "invalid",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refresh.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "grant_type is unsupported",
			expErr:  "grant_type is unsupported",
		}, {
			name: "Invalid client_id",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     "invalid",
				"client_secret": s.config().CAdmin.Secret,
				"refresh_token": refresh.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "client not found",
			expErr:  "client not found",
		}, {
			name: "Invalid client_secret",
			query: map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     s.config().CAdmin.Id,
				"client_secret": "invalid",
				"refresh_token": refresh.Hash,
			},
			expCode: http.StatusBadRequest,
			expBody: "client not found",
			expErr:  "client not found",
		},
	}

	ctrl := oauth.NewTokenController(s.app.Provider.OAuth())

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			query := s.buildQuery(tc.query)

			req := httptest.NewRequest(http.MethodPost, "/?"+query, nil)
			req.Header.Add("Content-Type", echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := s.app.HttpServer.NewContext(req, rec)

			if err = s.sendToServer(ctrl.Token, c); err != nil {
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
