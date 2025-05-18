package main

type Receiver interface {
	Receive(ctx Context)
}

type ReceiverFunc func(ctx Context)

func (f ReceiverFunc) Receive(ctx Context) { f(ctx) }
