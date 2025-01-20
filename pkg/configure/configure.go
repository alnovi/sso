package configure

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

func ParseEnv(config any) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := envconfig.Process(ctx, config)
	if err != nil {
		return fmt.Errorf("envconfig: %w", err)
	}

	return nil
}
