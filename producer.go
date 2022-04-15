package queue

type Producer[T any] interface {
	Produce(value T) error
}
