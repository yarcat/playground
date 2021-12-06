package vase

import (
	"log"
	"time"
)

type hammer struct{ v *Vase }

func Break(v *Vase) hammer { return hammer{v} }

func (h hammer) Every(d time.Duration) (cancel func()) {
	done := make(chan struct{}, 1)
	go func() {
		t := time.NewTicker(d)
		for {
			select {
			case <-done:
				return
			case <-t.C:
				if err := h.v.Close(); err != nil {
					log.Printf("(ignored) failed to break %s: %v", h.v.dir, err)
				}
			}
		}
	}()
	return func() { done <- struct{}{} }
}
