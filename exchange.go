package queue

import "context"

type Exchange[T any, U any] interface {
	Producer[T]
	Consumer[U]
}

var _ Exchange[any, any] = (*exchange[any, any])(nil)

type exchange[T any, U any] struct {
	producer Producer[T]
	consumer Consumer[U]
}

func (e *exchange[T, U]) Produce(value T) error {
	return e.producer.Produce(value)
}

func (e *exchange[T, U]) Consume(ctx context.Context) (U, error) {
	return e.consumer.Consume(ctx)
}

func (e *exchange[T, U]) Cancel() {
	e.consumer.Cancel()
}

func NewExchange[T any, U any](producerQueue Queue[T], consumerQueue Queue[U]) (Exchange[T, U], error) {
	consumer, err := consumerQueue.NewConsumer()
	if err != nil {
		return nil, err
	}
	return &exchange[T, U]{
		producer: producerQueue,
		consumer: consumer,
	}, nil
}
