package queue

import (
	"context"
	"fmt"
	"github.com/raf924/queue/utils"
)

type Consumer[T any] interface {
	Consume(context.Context) (T, error)
	Cancel()
}

var _ Consumer[any] = (*consumer[any])(nil)

type consumer[T any] struct {
	id     string
	buffer *linkedBuffer[T]
	q      cancellable
}

func (c *consumer[T]) Consume(ctx context.Context) (T, error) {
	if c.buffer.cancelled {
		return *new(T), fmt.Errorf("consumer was cancelled")
	}
	buffer := c.buffer
	value, err := utils.SyncGenerator[T, utils.GeneratorFuncWithError[T]](buffer.rwm, ctx, func(ctx context.Context) (T, error) {
		if buffer.root == nil {
			waitContext, cancel := context.WithCancel(ctx)
			select {
			case <-ctx.Done():
				buffer.c.Signal()
				buffer.rwm.Lock()
				return *new(T), fmt.Errorf("consumer operation was cancelled")
			case <-utils.ChannelizeOnce[struct{}, utils.GeneratorFunc[struct{}]](waitContext, func(ctx context.Context) struct{} {
				buffer.c.Wait()
				return struct{}{}
			}, cancel):
			}
		}
		root := buffer.root
		buffer.root = buffer.root.next
		return root.value, nil
	})
	if err != nil {
		return *new(T), err
	}
	return value, nil
}

func (c *consumer[T]) Cancel() {
	c.q.cancel(c.id)
}

type cancellable interface {
	cancel(id string)
}
