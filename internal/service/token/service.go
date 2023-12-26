package token

import (
	"context"
	"errors"
	"time"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/pkg/rand"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ttlCode    = 1
	ttlAccess  = 30
	ttlRefresh = 14400
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) NewCode(ctx context.Context, user entity.User, client entity.Client, meta *entity.TokenMeta) (*entity.Token, error) {
	now := time.Now()

	token := &entity.Token{
		Class:      entity.TokenClassCode,
		Hash:       rand.Base62(40),
		UserId:     &user.Id,
		ClientId:   &client.Id,
		Meta:       meta,
		NotBefore:  now,
		Expiration: now.Add(time.Minute * ttlCode),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) NewAccess(_ context.Context, user entity.User, client entity.Client) (*entity.Token, error) {
	now := time.Now()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"client_id":    client.Id,
		"client_class": client.Class,
		"user_id":      user.Id,
		"user_image":   user.Image,
		"user_name":    user.Name,
		"user_email":   user.Email,
		"nbf":          now.Unix(),
		"exp":          now.Add(time.Minute * ttlAccess).Unix(),
	})

	hash, err := jwtToken.SignedString([]byte(client.Secret))
	if err != nil {
		return nil, err
	}

	token := &entity.Token{
		Class:      entity.TokenClassAccess,
		Hash:       hash,
		UserId:     &user.Id,
		ClientId:   &client.Id,
		NotBefore:  now,
		Expiration: now.Add(time.Minute * ttlAccess),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return token, nil
}

func (s *Service) NewRefresh(ctx context.Context, user entity.User, client entity.Client, meta *entity.TokenMeta) (*entity.Token, error) {
	now := time.Now()

	token := &entity.Token{
		Class:      entity.TokenClassRefresh,
		Hash:       rand.Base62(40),
		UserId:     &user.Id,
		ClientId:   &client.Id,
		Meta:       meta,
		NotBefore:  now,
		Expiration: now.Add(time.Minute * ttlRefresh),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.CreateToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) FindToken(ctx context.Context, client entity.Client, class, hash string) (*entity.Token, error) {
	token, err := s.repo.GetTokenByClassAndHash(ctx, class, hash)
	if errors.Is(err, exception.TokenNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, exception.Wrap(exception.TokenNotFound, err)
	}

	if *token.ClientId != client.Id {
		return nil, exception.Wrap(exception.TokenNotFound, errors.New("token does not belong to the client"))
	}

	return token, err
}

func (s *Service) RemoveToken(ctx context.Context, tokenId string) error {
	return s.repo.DeleteToken(ctx, tokenId)
}
