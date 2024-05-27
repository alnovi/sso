package secure

import (
	"context"
	"time"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/rand"
	"github.com/golang-jwt/jwt/v4"
)

const (
	ttlCode    = 1
	ttlAccess  = 5
	ttlRefresh = 60
	ttlReset   = 60

	costCode    = 50
	costRefresh = 50
	costReset   = 50
)

type repository interface {
	CreateToken(ctx context.Context, token *entity.Token) error
}

type Secure struct {
	repo repository
}

type AccessClaims struct {
	jwt.RegisteredClaims
	UserName  string `json:"name"`
	UserEmail string `json:"email"`
	UserRole  string `json:"role"`
}

func NewSecure(repo repository) *Secure {
	return &Secure{repo: repo}
}

func (s *Secure) NewCodeToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error) {
	token := &entity.Token{
		Class:      entity.TokenClassCode,
		Hash:       rand.Base62(costCode),
		UserID:     user.ID,
		ClientID:   client.ID,
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(ttlCode * time.Minute),
	}

	return token, s.repo.CreateToken(ctx, token)
}

func (s *Secure) NewAccessToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error) {
	now := time.Now()
	expiration := now.Add(ttlAccess * time.Minute)

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
		UserName:  user.Name,
		UserEmail: user.Email,
		UserRole:  "user",
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hash, err := jwtToken.SignedString([]byte(client.Secret))
	if err != nil {
		return nil, err
	}

	token := &entity.Token{
		Class:      entity.TokenClassAccess,
		Hash:       hash,
		UserID:     user.ID,
		ClientID:   client.ID,
		NotBefore:  now,
		Expiration: expiration,
	}

	return token, s.repo.CreateToken(ctx, token)
}

func (s *Secure) NewRefreshToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error) {
	now := time.Now()

	token := &entity.Token{
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(costRefresh),
		UserID:     user.ID,
		ClientID:   client.ID,
		NotBefore:  now,
		Expiration: now.Add(ttlRefresh * time.Minute),
	}

	return token, s.repo.CreateToken(ctx, token)
}
