package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (j *JWT) GenerateAccessToken(clientId, userId, role string) (AccessClaims, string, error) {
	now := time.Now()
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    clientId,
			Subject:   userId,
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
		Role: role,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		return claims, "", fmt.Errorf("could not sign jwt access token: %w", err)
	}

	return claims, token, nil
}
