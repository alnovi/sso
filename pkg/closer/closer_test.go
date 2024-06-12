package closer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloser_Success(t *testing.T) {
	count := 5
	completed := 0

	c := New()

	for i := 0; i < count; i++ {
		c.Add(func(ctx context.Context) error {
			completed++
			return nil
		})
	}

	err := c.Close(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, count, completed)
}

func TestCloser_Error(t *testing.T) {
	count := 5
	completed := 0

	c := New()

	for i := 0; i < count; i++ {
		c.Add(func(ctx context.Context) error {
			completed++
			return fmt.Errorf("fail task")
		})
	}

	err := c.Close(context.Background())
	assert.Error(t, err)
	assert.Equal(t, count, completed)
}

func TestCloser_Timeout(t *testing.T) {
	count := 5
	completed := 0

	c := New()

	for i := 0; i < count; i++ {
		c.Add(func(ctx context.Context) error {
			time.Sleep(time.Second)
			completed++
			return nil
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := c.Close(ctx)
	assert.Error(t, err)
	assert.NotEqual(t, count, completed)
}
