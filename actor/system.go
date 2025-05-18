package main

import (
	"log/slog"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type System struct {
	LoggerMixin
	ID              string
	ProcessRegistry *ProcessRegistry
	Root            Context
}

func InitSystem(logger *slog.Logger) *System {
	s := &System{
		LoggerMixin: InitLoggerMixin(logger),
		ID:          gonanoid.Must(),
	}
	s.ProcessRegistry = NewProcessRegistry(s)

	return s
}

func (s *System) LocalPID(id string) *PID { return &PID{ID: id} }

func (s *System) Send(pid *PID, msg any) {
	s.ProcessRegistry.FromID(pid.ID).SendMessage(msg)
}
