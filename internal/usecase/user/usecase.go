package user

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/service"
)

type UseCase struct {
	user service.User
}

func New(user service.User) *UseCase {
	return &UseCase{user: user}
}

func (uc *UseCase) UserInfo(ctx context.Context, userId string) (*entity.User, error) {
	return uc.user.GetUserById(ctx, userId)
}

func (uc *UseCase) UserUpdate(ctx context.Context, user *entity.User) error {
	return nil
}
