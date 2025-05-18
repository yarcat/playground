package main

import (
	"fmt"
	"log/slog"
	"time"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	logger := slog.Default()

	rootCtx := NewRootContext(InitSystem(logger))

	opts := &SpawnOptions{
		Receiver: ReceiverFunc(func(ctx Context) {
			fmt.Println(ctx.Self().ID, ctx.Message())
		}),
	}

	pid := rootCtx.Spawn(opts)
	rootCtx.Send(pid, "Hello, World!")
	rootCtx.Send(pid, "Hello, World!")
	rootCtx.Send(pid, "Hello, World!")
	rootCtx.Send(pid, "Hello, World!")

	time.Sleep(1 * time.Second)
	fmt.Println("Sending messages again...")

	rootCtx.Send(pid, "Hello, World!")
	rootCtx.Send(pid, "Hello, World!")
	rootCtx.Send(pid, "Hello, World!")
	rootCtx.Send(pid, "Hello, World!")

	time.Sleep(1 * time.Second)
}
