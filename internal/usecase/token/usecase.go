package token

import (
	"context"
	"fmt"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/service"
)

const (
	grantTypeAuthorizationCode = "authorization_code"
	grantTypeRefreshToken      = "refresh_token"
)

type UseCase struct {
	token service.Token
	user  service.User
}

func New(token service.Token, user service.User) *UseCase {
	return &UseCase{token: token, user: user}
}

func (uc *UseCase) AccessAndRefreshToken(ctx context.Context, dto dto.AccessToken) (*entity.TokenWithUser, *entity.Token, error) {
	var class string
	var hash string

	switch dto.GrantType {
	case grantTypeAuthorizationCode:
		class = entity.TokenClassCode
		hash = dto.Code
	case grantTypeRefreshToken:
		class = entity.TokenClassRefresh
		hash = dto.Refresh
	default:
		return nil, nil, exception.Wrap(exception.TokenNotFound, fmt.Errorf("grant_type '%s' is not supported", dto.GrantType))
	}

	token, err := uc.token.FindToken(ctx, dto.Client, class, hash)
	if err != nil {
		return nil, nil, err
	}

	_ = uc.token.RemoveToken(ctx, token.Id)

	user, err := uc.user.GetUserById(ctx, *token.UserId)
	if err != nil {
		return nil, nil, exception.Wrap(exception.TokenNotFound, err)
	}

	if err = uc.user.CanUseClient(ctx, *user, dto.Client); err != nil {
		return nil, nil, err
	}

	access, err := uc.token.NewAccess(ctx, *user, dto.Client)
	if err != nil {
		return nil, nil, err
	}

	refresh, err := uc.token.NewRefresh(ctx, *user, dto.Client, token.Meta)
	if err != nil {
		return nil, nil, err
	}

	accessWithUser := &entity.TokenWithUser{
		Token: *access,
		User:  *user,
	}

	return accessWithUser, refresh, nil
}
