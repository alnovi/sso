package token

import (
	"context"
	"time"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service/jwt"
	"github.com/alnovi/sso/pkg/rand"
)

const (
	ClassCode    = "code"
	ClassAccess  = "access"
	ClassRefresh = "refresh"
	ClassForgot  = "forgot"
	costCode     = 50
	costRefresh  = 50
	costForgot   = 50
)

type Token struct {
	repo repository.Repository
	jwt  *jwt.JWT
}

func New(repo repository.Repository, jwt *jwt.JWT) *Token {
	return &Token{repo: repo, jwt: jwt}
}

func (s *Token) CodeToken(ctx context.Context, sessionId, clientId, userId string) (*entity.Token, error) {
	token := &entity.Token{
		Class:      ClassCode,
		Hash:       rand.Base62(costCode),
		SessionId:  &sessionId,
		UserId:     &userId,
		ClientId:   &clientId,
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(time.Minute),
	}

	err := s.repo.TokenCreate(ctx, token)

	return token, err
}

func (s *Token) AccessToken(_ context.Context, sessionId, clientId, userId, role string) (*entity.Token, error) {
	claims, hash, err := s.jwt.GenerateAccessToken(clientId, userId, role)
	if err != nil {
		return nil, err
	}

	token := &entity.Token{
		Class:      ClassAccess,
		Hash:       hash,
		SessionId:  &sessionId,
		UserId:     &userId,
		ClientId:   &clientId,
		NotBefore:  claims.NotBefore(),
		Expiration: claims.ExpiresAt(),
	}

	return token, nil
}

func (s *Token) RefreshToken(ctx context.Context, sessionId, clientId, userId string, before time.Time) (*entity.Token, error) {
	token := &entity.Token{
		Class:      ClassRefresh,
		Hash:       rand.Base62(costRefresh),
		SessionId:  &sessionId,
		UserId:     &userId,
		ClientId:   &clientId,
		NotBefore:  before,
		Expiration: before.AddDate(0, 1, 0),
	}

	err := s.repo.TokenCreate(ctx, token)

	return token, err
}

func (s *Token) ForgotPasswordToken(ctx context.Context, clientId, userId, query string) (*entity.Token, error) {
	token := &entity.Token{
		Class:      ClassForgot,
		Hash:       rand.Base62(costForgot),
		UserId:     &userId,
		ClientId:   &clientId,
		Payload:    entity.Payload{entity.PayloadQuery: query},
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(time.Hour),
	}

	err := s.repo.TokenCreate(ctx, token)

	return token, err
}
