package token

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
	"github.com/alnovi/sso/pkg/rand"
	"github.com/alnovi/sso/pkg/utils"
)

var (
	ErrInvalidPrivateKey = errors.New("invalid rsa private key")
	ErrInvalidPublicKey  = errors.New("invalid rsa public key")
	ErrTokenNotFound     = errors.New("token not found")
)

type Token struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	repo       *repository.Repository
}

func New(prvKey *rsa.PrivateKey, pubKey *rsa.PublicKey, repo *repository.Repository) (*Token, error) {
	if prvKey == nil {
		return nil, fmt.Errorf("%w: key is nil", ErrInvalidPrivateKey)
	}

	if pubKey == nil {
		return nil, fmt.Errorf("%w: key is nil", ErrInvalidPublicKey)
	}

	return &Token{
		privateKey: prvKey,
		publicKey:  pubKey,
		repo:       repo,
	}, nil
}

func (t *Token) CodeToken(ctx context.Context, sessionId, clientId, userId string) (*entity.Token, error) {
	ctx, span := helper.SpanStart(ctx, "Token.CodeToken")
	defer span.End()

	token := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassCode,
		Hash:       rand.Base62(entity.TokenCodeCost),
		SessionId:  utils.Point(sessionId),
		UserId:     utils.Point(userId),
		ClientId:   utils.Point(clientId),
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(entity.TokenCodeTTL),
	}

	if err := t.repo.TokenCreate(ctx, token); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return token, nil
}

func (t *Token) AccessToken(ctx context.Context, sessionId, clientId, userId, userName, role string, opts ...Option) (*entity.Token, error) {
	_, span := helper.SpanStart(ctx, "Token.AccessToken")
	defer span.End()

	now := time.Now()

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(entity.TokenAccessTTL)),
		},
		Session: sessionId,
		Client:  clientId,
		User:    userId,
		Name:    userName,
		Role:    role,
	}

	t.applyOptions(&claims, opts)

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(t.privateKey)
	if err != nil {
		helper.SpanError(span, fmt.Errorf("could not sign jwt access token: %w", err))
		return nil, fmt.Errorf("could not sign jwt access token: %w", err)
	}

	access := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassAccess,
		Hash:       token,
		SessionId:  utils.Point(sessionId),
		UserId:     utils.Point(userId),
		ClientId:   utils.Point(clientId),
		NotBefore:  claims.NotBefore(),
		Expiration: claims.ExpiresAt(),
	}

	return access, nil
}

func (t *Token) ValidateAccessToken(_ context.Context, access string) (*AccessClaims, error) {
	_, span := helper.SpanStart(context.Background(), "Token.ValidateAccessToken")
	defer span.End()

	access = strings.TrimPrefix(access, "Bearer ")
	access = strings.TrimSpace(access)

	token, err := jwt.ParseWithClaims(access, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.publicKey, nil
	})

	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	if token == nil {
		helper.SpanError(span, errors.New("invalid token"))
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		helper.SpanError(span, errors.New("invalid token"))
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (t *Token) RefreshToken(ctx context.Context, sessionId, clientId, userId string, notBefore time.Time, opts ...Option) (*entity.Token, error) {
	ctx, span := helper.SpanStart(ctx, "Token.RefreshToken")
	defer span.End()

	refresh := &entity.Token{
		Id:         uuid.NewString(),
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(entity.TokenRefreshCost),
		SessionId:  utils.Point(sessionId),
		ClientId:   utils.Point(clientId),
		UserId:     utils.Point(userId),
		NotBefore:  notBefore,
		Expiration: notBefore.Add(entity.TokenRefreshTTL),
	}

	t.applyOptions(refresh, opts)

	if err := t.repo.TokenCreate(ctx, refresh); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return refresh, nil
}

func (t *Token) ValidateRefreshToken(ctx context.Context, refresh string) (*entity.Token, error) {
	ctx, span := helper.SpanStart(ctx, "Token.ValidateRefreshToken")
	defer span.End()

	refresh = strings.TrimSpace(refresh)

	token, err := t.repo.TokenByHash(ctx, refresh, repository.Class(entity.TokenClassRefresh))
	if err != nil {
		helper.SpanError(span, fmt.Errorf("%w: %s", ErrTokenNotFound, err))
		return nil, fmt.Errorf("%w: %s", ErrTokenNotFound, err)
	}

	if !token.IsActive() {
		helper.SpanError(span, fmt.Errorf("%w: tiken is inactive", ErrTokenNotFound))
		return nil, fmt.Errorf("%w: tiken is inactive", ErrTokenNotFound)
	}

	return token, nil
}

func (t *Token) ForgotPasswordToken(ctx context.Context, clientId, userId, query, ip, agent string, opts ...Option) (*entity.Token, error) {
	ctx, span := helper.SpanStart(ctx, "Token.ForgotPasswordToken")
	defer span.End()

	forgot := &entity.Token{
		Id:       uuid.NewString(),
		Class:    entity.TokenClassForgot,
		Hash:     rand.Base62(entity.TokenForgotCost),
		ClientId: utils.Point(clientId),
		UserId:   utils.Point(userId),
		Payload: entity.Payload{
			entity.PayloadQuery: query,
			entity.PayloadIP:    ip,
			entity.PayloadAgent: agent,
		},
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(time.Hour),
	}

	t.applyOptions(forgot, opts)

	if err := t.repo.TokenCreate(ctx, forgot); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return forgot, nil
}

func (t *Token) ValidateForgotToken(ctx context.Context, forgot string) (*entity.Token, error) {
	ctx, span := helper.SpanStart(ctx, "Token.ValidateForgotToken")
	defer span.End()

	forgot = strings.TrimSpace(forgot)

	token, err := t.repo.TokenByHash(ctx, forgot, repository.Class(entity.TokenClassForgot))
	if err != nil {
		helper.SpanError(span, fmt.Errorf("%w: %s", ErrTokenNotFound, err))
		return nil, fmt.Errorf("%w: %s", ErrTokenNotFound, err)
	}

	if !token.IsActive() {
		helper.SpanError(span, fmt.Errorf("%w: tiken is inactive", ErrTokenNotFound))
		return nil, fmt.Errorf("%w: tiken is inactive", ErrTokenNotFound)
	}

	return token, nil
}

func (t *Token) applyOptions(e any, opts []Option) {
	for _, opt := range opts {
		opt(e)
	}
}
