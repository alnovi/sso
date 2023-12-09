package configure

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

func LoadFromEnv[C any](config C) (C, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := envconfig.Process(ctx, &config)
	if err != nil {
		return config, fmt.Errorf("failed to envconfig.Process: %w", err)
	}

	return config, nil
}
