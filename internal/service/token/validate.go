package token

import (
	"context"
	"errors"

	"github.com/alnovi/sso/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) Validate(ctx context.Context, clientId, secret, class, hash string) (*entity.Token, error) {
	if class == entity.TokenClassAccess {
		return s.validateJWT(clientId, secret, hash)
	}

	return s.validateToken(ctx, clientId, class, hash)
}

func (s *Service) validateJWT(clientId, secret, hash string) (*entity.Token, error) {
	jwtToken, err := jwt.Parse(hash, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("can't parse jwt token")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if jwtToken == nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("can't parse claims")
	}

	userId := claims["user_id"].(string)

	if clientId != claims["client_id"].(string) {
		return nil, errors.New("invalid token")
	}

	notBefore, err := claims.GetNotBefore()
	if err != nil {
		return nil, err
	}

	expiration, err := claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	token := &entity.Token{
		Class:      entity.TokenClassAccess,
		Hash:       hash,
		UserId:     &userId,
		ClientId:   &clientId,
		NotBefore:  notBefore.Time,
		Expiration: expiration.Time,
		CreatedAt:  notBefore.Time,
		UpdatedAt:  notBefore.Time,
	}

	return token, nil
}

func (s *Service) validateToken(ctx context.Context, clientId, class, hash string) (*entity.Token, error) {
	token, err := s.repo.GetTokenByClassAndHash(ctx, class, hash)
	if err != nil {
		return nil, err
	}

	if token.ClientId == nil || clientId != *token.ClientId {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
