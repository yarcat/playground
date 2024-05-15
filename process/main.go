package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

type Worker struct {
	cmd *exec.Cmd
}

func NewWorker(ctx context.Context) (*Worker, error) {
	p, err := exec.LookPath("python3")
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, p, "task.py")
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Worker{cmd: cmd}, nil
}

func (w *Worker) Wait() error { return w.cmd.Wait() }

func runPool(ctx context.Context, n int) {
	var wg sync.WaitGroup
	defer wg.Wait()

	workers := make(chan struct{}, n)
	for {
		select {
		case <-ctx.Done():
			return
		case workers <- struct{}{}:
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() { <-workers }()
				w, err := NewWorker(ctx)
				if err == nil {
					err = w.Wait()
				}
				fmt.Println(err)
				var e *exec.ExitError
				if errors.As(err, &e) {
					fmt.Println(e.Sys().(syscall.WaitStatus).Signaled())
				}
			}()
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer time.AfterFunc(5*time.Second, cancel).Stop()
	runPool(ctx, 5)
}
