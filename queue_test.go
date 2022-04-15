package queue

import (
	"context"
	"testing"
	"time"
)

func TestConsume(t *testing.T) {
	queue := NewQueue[int]()
	consumer, err := queue.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}
	err = queue.Produce(5)
	if err != nil {
		t.Fatal(err)
	}
	val, err := consumer.Consume(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if val != 5 {
		t.Fatal("expected 5 got", val)
	}
}

func TestSelect(t *testing.T) {
	queue := NewQueue[int]()
	consumer, err := queue.NewConsumer()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	ch := make(chan int)
	go func() {
		for ctx.Err() == nil {
			c, err := consumer.Consume(ctx)
			if err != nil {
				cancel()
			}
			ch <- c
		}
	}()
	select {
	case <-ctx.Done():
	case <-ch:
		cancel()
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatal("expected", context.DeadlineExceeded, "got", ctx.Err())
	}
}
