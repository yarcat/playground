package main

type Scheduler interface {
	Schedule(func())
}

type SchedulerFunc func(func())

func (f SchedulerFunc) Schedule(fx func()) { f(fx) }
