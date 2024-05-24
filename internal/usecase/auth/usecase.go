package auth

import (
	"context"
	"fmt"
	"slices"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/pkg/rand"
)

const (
	costTokenCode = 50
	ttlTokenCode  = 1
)

type repository interface {
	ClientByID(ctx context.Context, id string) (*entity.Client, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateToken(ctx context.Context, token *entity.Token) error
}

type UseCase struct {
	repo repository
}

func NewUseCase(repo repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) ValidateGrantType(ctx context.Context, inp dto.InpValidateGrantType) (*entity.Client, error) {
	if inp.ClientID == "" {
		return nil, exception.ErrClientNotFound
	}

	if inp.GrantType == "" {
		return nil, exception.ErrUnsupportedGrantType
	}

	client, err := uc.repo.ClientByID(ctx, inp.ClientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrClientNotFound, err)
	}

	if !client.IsActive {
		return nil, exception.ErrClientNotFound
	}

	if !slices.Contains(client.GrantTypes, inp.GrantType) {
		return nil, exception.ErrUnsupportedGrantType
	}

	if inp.RedirectURI != "" && inp.RedirectURI != client.Callback {
		return nil, exception.ErrClientNotFound
	}

	return client, nil
}

func (uc *UseCase) CodeByCredentials(ctx context.Context, inp dto.InpAuthByCredentials) (*entity.Token, error) {
	user, err := uc.repo.UserByEmail(ctx, inp.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrUserNotFound, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inp.Password)); err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrPasswordIncorrect, err)
	}

	token := &entity.Token{
		Type:       entity.TokenTypeCode,
		Hash:       rand.Base62(costTokenCode),
		UserID:     user.ID,
		ClientID:   inp.Client.ID,
		NotBefore:  time.Now(),
		Expiration: time.Now().Add(ttlTokenCode * time.Minute),
	}

	if err = uc.repo.CreateToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}
