package main

type Context interface {
	Spawn(*SpawnOptions) *PID
	Self() *PID
	Message() any
	Send(*PID, any)
}

type rootContext struct {
	s *System
}

func NewRootContext(s *System) *rootContext {
	return &rootContext{s: s}
}

func (rc *rootContext) Self() *PID                    { return nil }
func (rc *rootContext) Message() any                  { return nil }
func (rc *rootContext) Spawn(opts *SpawnOptions) *PID { return opts.Spawn(rc.s, rc) }
func (rc *rootContext) Send(pid *PID, msg any)        { rc.s.Send(pid, msg) }
