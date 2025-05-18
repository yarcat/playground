package main

import (
	"sync"
	"sync/atomic"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type PID struct {
	ID string
}

type PIDBuilder interface {
	LocalPID(id string) *PID
}

type Process interface {
	SendMessage(any)
}

type ProcessRegistry struct {
	processes sync.Map // ID -> Process.
	pid       PIDBuilder
}

func NewProcessRegistry(pid PIDBuilder) *ProcessRegistry {
	return &ProcessRegistry{
		pid:       pid,
		processes: sync.Map{},
	}
}

func (*ProcessRegistry) NextID() string { return gonanoid.Must() }

func (pr *ProcessRegistry) FromID(id string) Process {
	if v, ok := pr.processes.Load(id); ok {
		return v.(Process)
	}
	return nil
}

func (pr *ProcessRegistry) AddProcess(id string, process Process) *PID {
	pr.processes.Store(id, process)
	return pr.pid.LocalPID(id)
}

type process struct {
	Queue
	Receiver
	Scheduler
	Context
	*PID
	SpawnOptions *SpawnOptions

	len     atomic.Int32
	running atomic.Bool
}

func (p *process) SendMessage(msg any) {
	p.Push(msg)
	p.len.Add(1)
	if p.running.CompareAndSwap(false, true) {
		p.Schedule(p.process)
	}
}

func (p *process) process() {
	defer p.running.Store(false)

	for {
		msg := p.Pop()
		if msg == nil {
			break
		}

		p.len.Add(-1)
		p.Receiver.Receive(&processContext{p, msg})
	}
}

type processContext struct {
	*process
	msg any
}

func (p *processContext) Message() any { return p.msg }
func (p *processContext) Self() *PID   { return p.PID }
