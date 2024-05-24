package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Func func(ctx context.Context) error

type Closer struct {
	mu sync.Mutex
	fs []Func
}

func New() *Closer {
	return &Closer{}
}

func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fs = append(c.fs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	messages := make([]string, 0, len(c.fs))
	complete := make(chan struct{}, 1)

	go func() {
		for i := len(c.fs) - 1; i >= 0; i-- {
			if err := c.fs[i](ctx); err != nil {
				messages = append(messages, fmt.Sprintf("[!] %s", err.Error()))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %s", ctx.Err())
	}

	if len(messages) > 0 {
		return fmt.Errorf("shutdown finished with error(s): \n%s", strings.Join(messages, "\n"))
	}

	return nil
}
