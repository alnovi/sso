package sessions

import (
	"context"
	"errors"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

var (
	ErrSessionNotFound = errors.New("session not found")
)

type Session struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Session {
	return &Session{repo: repo}
}

func (s *Session) Create(ctx context.Context, userId, ip, agent string) (*entity.Session, error) {
	session := &entity.Session{
		UserId: userId,
		Ip:     ip,
		Agent:  agent,
	}

	err := s.repo.SessionCreate(ctx, session)

	return session, err
}

func (s *Session) Delete(ctx context.Context, sessionId string) error {
	err := s.repo.SessionDelete(ctx, sessionId)
	if err != nil && !errors.Is(err, repository.ErrNoResults) {
		return err
	}
	return nil
}

func (s *Session) GetById(ctx context.Context, sessionId string) (*entity.Session, error) {
	session, err := s.repo.SessionById(ctx, sessionId)
	if errors.Is(err, repository.ErrNoResults) {
		return nil, ErrSessionNotFound
	}
	return session, err
}
