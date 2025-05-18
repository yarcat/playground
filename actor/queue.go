package main

type Queue interface {
	Push(any)
	Pop() any
}

type QueueFactory interface {
	Queue() Queue
}

type QueueFactoryFunc func() Queue

func (f QueueFactoryFunc) Queue() Queue { return f() }

type queue chan any

func NewQueue(size int) queue {
	return make(queue, size)
}

func (q queue) Push(msg any) { q <- msg }
func (q queue) Pop() any {
	select {
	case msg := <-q:
		return msg
	default:
		return nil
	}
}
