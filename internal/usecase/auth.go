package usecase

import (
	"context"
	"fmt"
	"slices"

	"golang.org/x/crypto/bcrypt"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

func (uc *UseCase) ValidateResponseType(ctx context.Context, inp dto.ValidateResponseType) (*entity.Client, error) {
	if inp.ClientID == "" {
		return nil, exception.ErrClientNotFound
	}

	if !slices.Contains([]string{dto.ResponseTypeCode}, inp.ResponseType) {
		return nil, exception.ErrUnsupportedGrantType
	}

	client, err := uc.repo.ClientByID(ctx, inp.ClientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrClientNotFound, err)
	}

	if !client.IsActive {
		return nil, exception.ErrClientNotFound
	}

	return client, nil
}

func (uc *UseCase) CodeByCredentials(ctx context.Context, inp dto.AuthByCredentials) (*entity.Token, error) {
	user, err := uc.repo.UserByEmail(ctx, inp.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrUserNotFound, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inp.Password)); err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrPasswordIncorrect, err)
	}

	return uc.secure.NewCodeToken(ctx, user, inp.Client)
}

func (uc *UseCase) AccessTokenByCode(ctx context.Context, inp dto.AccessTokenByCode) (*entity.Token, *entity.Token, error) {
	client, err := uc.repo.ClientByIdAndSecret(ctx, inp.ClientID, inp.ClientSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", exception.ErrClientNotFound, err)
	}

	token, err := uc.repo.ClientTokenByClassAndHash(ctx, client.ID, entity.TokenClassCode, inp.CodeHash)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", exception.ErrTokenNotFound, err)
	}

	if err = uc.repo.DeleteTokenById(ctx, token.ID); err != nil {
		return nil, nil, fmt.Errorf("delete code token: %s", err)
	}

	user, err := uc.repo.UserByID(ctx, token.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", exception.ErrUserNotFound, err)
	}

	access, err := uc.secure.NewAccessToken(ctx, user, client)
	if err != nil {
		return nil, nil, err
	}

	refresh, err := uc.secure.NewRefreshToken(ctx, user, client, inp.IP, inp.Agent)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}

func (uc *UseCase) AccessTokenByRefresh(ctx context.Context, inp dto.AccessTokenByRefresh) (*entity.Token, *entity.Token, error) {
	client, err := uc.repo.ClientByIdAndSecret(ctx, inp.ClientID, inp.ClientSecret)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", exception.ErrClientNotFound, err)
	}

	token, err := uc.repo.ClientTokenByClassAndHash(ctx, client.ID, entity.TokenClassRefresh, inp.RefreshHash)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", exception.ErrTokenNotFound, err)
	}

	if err = uc.repo.DeleteTokenById(ctx, token.ID); err != nil {
		return nil, nil, fmt.Errorf("delete code token: %s", err)
	}

	user, err := uc.repo.UserByID(ctx, token.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", exception.ErrUserNotFound, err)
	}

	access, err := uc.secure.NewAccessToken(ctx, user, client)
	if err != nil {
		return nil, nil, err
	}

	refresh, err := uc.secure.NewRefreshToken(ctx, user, client, inp.IP, inp.Agent)
	if err != nil {
		return nil, nil, err
	}

	return access, refresh, nil
}
