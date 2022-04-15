package queue

import "sync"

type linkedBuffer[T any] struct {
	root      *bufferValue[T]
	rwm       *sync.RWMutex
	c         *sync.Cond
	cancelled bool
}

type bufferValue[T any] struct {
	value T
	next  *bufferValue[T]
}

func (buffer *linkedBuffer[T]) pushValue(value T) {
	qv := &bufferValue[T]{
		value: value,
	}
	buffer.rwm.RLock()
	root := buffer.root
	buffer.rwm.RUnlock()
	if root == nil {
		buffer.rwm.Lock()
		buffer.root = qv
		buffer.rwm.Unlock()
	} else {
		buffer.rwm.RLock()
		for root.next != nil {
			root = root.next
		}
		buffer.rwm.RUnlock()
		buffer.rwm.Lock()
		root.next = qv
		buffer.rwm.Unlock()
	}
	buffer.c.Signal()
}
