package configure

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

func ParseEnv(ctx context.Context, config any) error {
	return envconfig.Process(ctx, config)
}
