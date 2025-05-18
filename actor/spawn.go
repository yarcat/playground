package main

type SpawnOptions struct {
	Receiver  Receiver
	queueFact QueueFactory
	spawner   Spawner
	scheduler Scheduler
}

type Spawner interface {
	Spawn(s *System) *PID
}

func (so *SpawnOptions) Spawn(s *System, c Context) *PID {
	if so.spawner != nil {
		return so.spawner.Spawn(s)
	}
	return defaultSpawner(s, c, so)
}

func (so *SpawnOptions) getQueue() Queue {
	if so.queueFact != nil {
		return so.queueFact.Queue()
	}
	return NewQueue(10) // TODO: Make this configurable.
}

func (so *SpawnOptions) getScheduler() Scheduler {
	if so.scheduler != nil {
		return so.scheduler
	}
	return SchedulerFunc(func(f func()) { go f() })
}

func defaultSpawner(s *System, parentCtx Context, so *SpawnOptions) *PID {
	id := s.ProcessRegistry.NextID()

	p := &process{
		Queue:        so.getQueue(),
		Receiver:     so.Receiver,
		Scheduler:    so.getScheduler(),
		Context:      parentCtx,
		SpawnOptions: so,
	}
	pid := s.ProcessRegistry.AddProcess(id, p)
	p.PID = pid

	return pid
}
