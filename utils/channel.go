package utils

import "context"

type Generator[T any] interface {
	Generate(ctx context.Context) (T, error)
}

type generator[T any] interface {
	Generate(ctx context.Context) (T, error)
	GeneratorFunc[T] | GeneratorFuncWithError[T]
}

var _ Generator[any] = (GeneratorFunc[any])(nil)
var _ Generator[any] = (GeneratorFuncWithError[any])(nil)

type GeneratorFunc[T any] func(context.Context) T

func (f GeneratorFunc[T]) Generate(ctx context.Context) (T, error) {
	return f(ctx), nil
}

type GeneratorFuncWithError[T any] func(context.Context) (T, error)

func (f GeneratorFuncWithError[T]) Generate(ctx context.Context) (T, error) {
	return f(ctx)
}

func Channelize[T any, G Generator[T]](ctx context.Context, generator G, cancel context.CancelFunc) <-chan T {
	ch := make(chan T)
	go func() {
		for ctx.Err() == nil {
			generated, err := generator.Generate(ctx)
			if err != nil {
				cancel()
			}
			ch <- generated
		}
	}()
	return ch
}

func ChannelizeOnce[T any, G Generator[T]](ctx context.Context, generator G, cancel context.CancelFunc) <-chan T {
	ch := make(chan T)
	go func() {
		generated, err := generator.Generate(ctx)
		if err != nil {
			cancel()
		}
		ch <- generated
	}()
	return ch
}

func ContextChannelLoop[T any](ctx context.Context, channel <-chan T, handler func(T) error) {
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
		case value := <-channel:
			err := handler(value)
			if err != nil {
				return
			}
		}
	}
}
