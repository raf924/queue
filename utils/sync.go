package utils

import (
	"context"
	"sync"
)

func SyncGenerator[T any, G generator[T]](locker sync.Locker, ctx context.Context, generator G) (T, error) {
	locker.Lock()
	value, err := generator.Generate(ctx)
	locker.Unlock()
	return value, err
}
