package profile

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service"
)

type UseCase struct {
	client service.Client
	token  service.Token
	user   service.User
}

func New(client service.Client, token service.Token, user service.User) *UseCase {
	return &UseCase{client: client, token: token, user: user}
}

func (uc *UseCase) Profile(ctx context.Context, userId string) (*entity.User, []*entity.Token, []*entity.Client, error) {
	var err error

	user, err := uc.user.GetUserById(ctx, userId)
	if err != nil {
		return nil, nil, nil, err
	}

	tokens, err := uc.token.TokensByUserAndClass(ctx, userId, entity.TokenClassRefresh)
	if err != nil {
		return nil, nil, nil, err
	}

	//TODO: get access clients for user

	return user, tokens, nil, nil
}
