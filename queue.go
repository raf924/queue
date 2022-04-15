package queue

import (
	"github.com/segmentio/ksuid"
	"sync"
)

type Queue[T any] interface {
	Producer[T]
	NewConsumer() (Consumer[T], error)
}

var _ Queue[any] = (*queue[any])(nil)

func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		bufferMap: sync.Map{},
	}
}

type queue[T any] struct {
	bufferMap sync.Map
}

func (q *queue[T]) Produce(value T) error {
	q.bufferMap.Range(func(key, buffer any) bool {
		buffer.(*linkedBuffer[T]).pushValue(value)
		return true
	})
	return nil
}

func (q *queue[T]) cancel(id string) {
	if buffer, ok := q.bufferMap.Load(id); ok {
		buffer.(*linkedBuffer[T]).cancelled = true
		q.bufferMap.Delete(id)
	}
}

func (q *queue[T]) NewConsumer() (Consumer[T], error) {
	id := ksuid.New().String()
	rwm := new(sync.RWMutex)
	buffer := &linkedBuffer[T]{
		root: nil,
		rwm:  rwm,
		c:    sync.NewCond(rwm),
	}
	q.bufferMap.Store(id, buffer)
	return &consumer[T]{
		id:     id,
		buffer: buffer,
	}, nil
}
