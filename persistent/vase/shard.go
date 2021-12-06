package vase

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"path"
	"path/filepath"
	"sort"
	"time"
)

func Shards(v *Vase) ShardIterator {
	files, err := v.opts.Glob(path.Join(v.dir, "*"))
	if err != nil {
		log.Printf("(ignored) shards fail to glob %s: %v", v.dir, err)
	}
	sort.Strings(files)
	d := make(chan shardData, 10)
	done := make(chan struct{}, 1)
	go func() {
		defer close(d)
		if err := v.Flush(); err != nil {
			log.Printf("(ignored) failed to flush: %v", err)
		}
		for _, f := range files {
			r, err := v.opts.OpenFile(f)
			if err != nil {
				log.Printf("(ignored) unable to open shard: %v", err)
				continue
			}
			select {
			case d <- shardData{Name: f, R: &bufReadCloser{bufio.NewReader(r), r.Close}}:
			case <-done:
				return
			}
		}
	}()
	return ShardIterator{r: d, done: done}
}

type shardData struct {
	Name string
	R    io.ReadCloser
}

type ShardIterator struct {
	r    <-chan shardData
	done chan<- struct{}

	cur shardData
}

func (s *ShardIterator) Close() error {
	select {
	case s.done <- struct{}{}:
	default:
	}
	for d := range s.r {
		if err := d.R.Close(); err != nil {
			log.Printf("(ignored) error closing %s: %v", d.Name, err)
		}
	}
	return nil
}

func (s *ShardIterator) Next() bool {
	s.closeAndForget()
	d, ok := <-s.r
	if !ok {
		return false
	}
	s.cur = d
	return true
}

func (s *ShardIterator) closeAndForget() {
	if s.cur.R == nil {
		return
	}
	if err := s.cur.R.Close(); err != nil {
		log.Printf("(ignored) error closing %s: %v", s.cur.Name, err)
		s.cur.R, s.cur.Name = nil, ""
	}
}

func (s *ShardIterator) Read(f func(io.Reader) error) error {
	err := f(s.cur.R)
	if err == io.EOF {
		err = nil
	}
	return err
}

func (s *ShardIterator) Time() time.Time {
	var usec int64
	fmt.Sscanf(filepath.Base(s.cur.Name), "%x", &usec)
	return time.UnixMicro(usec)
}
