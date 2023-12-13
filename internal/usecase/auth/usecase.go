package auth

import (
	"context"
	"net/url"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service"
)

type UseCase struct {
	auth  service.Auth
	token service.Token
	user  service.User
}

func New(auth service.Auth, token service.Token, user service.User) *UseCase {
	return &UseCase{auth: auth, token: token, user: user}
}

func (uc *UseCase) AuthById(ctx context.Context, dto dto.AuthById) (*entity.User, *url.URL, error) {
	user, err := uc.user.GetUserById(ctx, dto.UserId)
	if err != nil {
		return nil, nil, err
	}

	err = uc.user.CanUseClient(ctx, *user, dto.Client)
	if err != nil {
		return nil, nil, err
	}

	callback, err := uc.buildCallback(ctx, *user, dto.Client, dto.IP, dto.Agent)
	if err != nil {
		return nil, nil, err
	}

	return user, callback, nil
}

func (uc *UseCase) AuthByCredentials(ctx context.Context, dto dto.AuthByCredentials) (*entity.User, *url.URL, error) {
	user, err := uc.auth.AuthByCredentials(ctx, dto.Login, dto.Password)
	if err != nil {
		return nil, nil, err
	}

	err = uc.user.CanUseClient(ctx, *user, dto.Client)
	if err != nil {
		return nil, nil, err
	}

	callback, err := uc.buildCallback(ctx, *user, dto.Client, dto.IP, dto.Agent)
	if err != nil {
		return nil, nil, err
	}

	return user, callback, nil
}

func (uc *UseCase) buildCallback(ctx context.Context, user entity.User, client entity.Client, ip, agent string) (*url.URL, error) {
	callback, err := url.Parse(client.Callback)
	if err != nil {
		return nil, err
	}

	meta := &entity.TokenMeta{
		IP:    ip,
		Agent: agent,
	}
	token, err := uc.token.NewCode(ctx, user, client, meta)
	if err != nil {
		return nil, err
	}

	query := callback.Query()
	query.Add("code", token.Hash)

	callback.RawQuery = query.Encode()

	return callback, nil
}
