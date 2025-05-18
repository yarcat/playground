package main

import "log/slog"

type LoggerMixin struct{ logger *slog.Logger }

func InitLoggerMixin(logger *slog.Logger) LoggerMixin {
	return LoggerMixin{logger: logger}
}

func (lm LoggerMixin) Logger() *slog.Logger { return lm.logger }
