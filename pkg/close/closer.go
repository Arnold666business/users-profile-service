package close

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
)

type Closer struct {
	mu      sync.Mutex
	closers []func(context.Context) error
}

func (c *Closer) Add(closers ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, closer := range closers {
		switch f := closer.(type) {
		case func(context.Context) error:
			c.closers = append(c.closers, f)
		case io.Closer:
			curCloser := f
			c.closers = append(c.closers, func(ctx context.Context) error {
				return curCloser.Close()
			})
		case func():
			c.closers = append(c.closers, func(ctx context.Context) error {
				f()
				return nil
			})
		default:
			log.Printf("unsupported close type: %T", closer)
		}
	}
}

func (c *Closer) CloseAllWithErrors(ctx context.Context) (closeErrors []error) {
	c.mu.Lock()
	closers := c.closers
	c.closers = nil
	c.mu.Unlock()

	for i := len(closers); i > 0; i-- {
		err := closers[i-1](ctx)
		if err != nil {
			closeErrors = append(closeErrors, err)
		}
	}
	return closeErrors
}

func (c *Closer) CloseAll(ctx context.Context) {
	errors := c.CloseAllWithErrors(ctx)
	if len(errors) > 0 {
		for _, err := range errors {
			log.Println("error closing close:", err)
		}
		os.Exit(1)
	}
}
