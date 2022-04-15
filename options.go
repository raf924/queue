package queue

type Option interface {
	Apply(q *queue[any])
}

type queueOptionFunc func(*queue[any])

var _ Option = queueOptionFunc(nil)

func (f queueOptionFunc) Apply(q *queue[any]) {
	f(q)
}
