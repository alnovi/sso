package token

import (
	"context"
	"time"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/rand"
)

const (
	classCode = "code"
	ttlCode   = 1
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) NewCode(ctx context.Context, userId, clientId, ip, agent string) (*entity.Token, error) {
	now := time.Now()
	meta := entity.TokenMeta{IP: ip, Agent: agent}

	token := &entity.Token{
		Class:      classCode,
		Hash:       rand.Base62(40),
		UserId:     &userId,
		ClientId:   &clientId,
		Meta:       &meta,
		NotBefore:  now,
		Expiration: now.Add(time.Minute * ttlCode),
	}

	if err := s.repo.CreateToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}
