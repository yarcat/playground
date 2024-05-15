package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

func runPool(ctx context.Context, n int) error {
	workers := make(chan struct{}, n)
	for {
		select {
		case <-ctx.Done():
			return nil
		case workers <- struct{}{}:
			go func() {
				defer func() { <-workers }()
				w, err := NewWorker(ctx)
				if err == nil {
					err = w.Wait()
				}
				fmt.Println(err)
			}()
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer time.AfterFunc(20*time.Second, cancel).Stop()
	fmt.Println(runPool(ctx, 5))
}
