package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/recoilme/sniper"
)

const expirationInterval = 24 * time.Hour

func main() {
	rand.Seed(time.Now().UnixNano())

	s, err := sniper.Open(
		sniper.Dir("timeseries"),
		sniper.SyncInterval(1*time.Second),
		sniper.ExpireInterval(expirationInterval),
	)
	if err != nil {
		log.Fatalf("failed to open sniper dir: %v", err)
	}
	defer s.Close()

	var wg sync.WaitGroup

	ctx, done := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		put(ctx, s)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		count(ctx, s)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	signal.Stop(c)
	log.Println("Waiting for workers to shutdown...")

	done()
	wg.Wait()
	log.Println("Everyone is done!")
}

func put(ctx context.Context, s *sniper.Store) {
	const interval = 10 * time.Millisecond
	t := time.NewTimer(interval)
	defer t.Stop()

	data := make([]TimeValue, 100)
	p := NewPutter(s)
	for {
		fill(data, time.Now())
		now := time.Now()
		err := p.Put(data)
		d := time.Since(now)
		if err == nil {
			log.Printf("INS: %v", d)
		} else {
			log.Printf("INS: %v", err)
		}

		t.Reset(interval)
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}
	}
}

func count(ctx context.Context, s *sniper.Store) {
	const interval = 10 * time.Second
	t := time.NewTimer(interval)
	defer t.Stop()
	for {
		now := time.Now()
		cnt := s.Count()
		d := time.Since(now)
		log.Printf("CNT: %d in %v", cnt, d)

		t.Reset(interval)
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}
	}
}

type TimeValue struct {
	Type int
	T    time.Time
	V    float32
}

func fill(tv []TimeValue, now time.Time) {
	for i := range tv {
		tv[i] = TimeValue{
			Type: i + 1,
			T:    now,
			V:    rand.Float32(),
		}
	}
}

type Putter struct {
	s *sniper.Store
}

func NewPutter(s *sniper.Store) *Putter { return &Putter{s: s} }

func (p *Putter) Put(tv []TimeValue) error {
	expire := time.Now().Add(expirationInterval)
	for i := range tv {
		k := p.KeyFor(&tv[i])
		switch err := p.s.Set(k, nil, uint32(expire.Unix())); err {
		case nil:
			continue
		case sniper.ErrCollision:
			log.Printf("INS: ignoring colliding entry: %v", &tv[i])
		default:
			log.Println(err)
			return err
		}
	}
	return nil
}

func (p *Putter) KeyFor(tv *TimeValue) []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, int32(tv.Type))
	binary.Write(&b, binary.BigEndian, tv.T.UnixMilli())
	binary.Write(&b, binary.BigEndian, math.Float32bits(tv.V))
	return b.Bytes()
}
