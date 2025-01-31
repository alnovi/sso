package mailing

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
)

type Mailing interface {
	Ping(ctx context.Context) error
	Close() error
	ForgotPassword(ctx context.Context, user *entity.User, token *entity.Token) error
}
