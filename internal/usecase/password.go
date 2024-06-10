package usecase

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

const (
	userPasswordCost = 10
)

func (uc *UseCase) ClientByResetPassword(ctx context.Context, hash string) (*entity.Client, error) {
	token, err := uc.repo.TokenByClassAndHash(ctx, entity.TokenClassResetPassword, hash)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrTokenNotFound, err)
	}

	if err = token.IsActive(); err != nil {
		return nil, err
	}

	client, err := uc.repo.ClientByID(ctx, token.ClientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrClientNotFound, err)
	}

	return client, nil
}

func (uc *UseCase) ForgotPassword(ctx context.Context, inp dto.ForgotPassword) (*entity.User, error) {
	user, err := uc.repo.UserByEmail(ctx, inp.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrUserNotFound, err)
	}

	token, err := uc.secure.NewResetPassword(ctx, user, inp.Client, inp.IP, inp.Agent)
	if err != nil {
		return nil, err
	}

	return user, uc.notify.ResetPassword(ctx, user, token)
}

func (uc *UseCase) ResetPassword(ctx context.Context, inp dto.ResetPassword) (*entity.Client, error) {
	token, err := uc.repo.TokenByClassAndHash(ctx, entity.TokenClassResetPassword, inp.Hash)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrTokenNotFound, err)
	}

	if err = token.IsActive(); err != nil {
		return nil, err
	}

	if err = uc.repo.DeleteTokenById(ctx, token.ID); err != nil {
		return nil, err
	}

	client, err := uc.repo.ClientByID(ctx, token.ClientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrClientNotFound, err)
	}

	user, err := uc.repo.UserByID(ctx, token.UserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", exception.ErrUserNotFound, err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(inp.Password), userPasswordCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashPassword)

	return client, uc.repo.UpdateUser(ctx, user)
}
