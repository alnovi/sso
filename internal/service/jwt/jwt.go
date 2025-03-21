package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/alnovi/sso/internal/entity"
)

var (
	ErrInvalidPrivateKey = errors.New("invalid rsa private key")
	ErrInvalidPublicKey  = errors.New("invalid rsa public key")
)

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New(prvKey, pubKey []byte) (*JWT, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(prvKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPrivateKey, err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPublicKey, err)
	}

	return &JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (j *JWT) AccessToken(sessionId, clientId, userId, role string) (AccessClaims, string, error) {
	now := time.Now()

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(entity.TokenAccessTTL)),
		},
		Session: sessionId,
		Client:  clientId,
		User:    userId,
		Role:    role,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		return claims, "", fmt.Errorf("could not sign jwt access token: %w", err)
	}

	return claims, token, nil
}

func (j *JWT) ParseToken(access string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(access, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
