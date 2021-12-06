package vase

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"sync"
)

type Vase struct {
	m    sync.Mutex
	w    io.WriteCloser
	dir  string
	opts options
}

func New(dir string, opts ...OptionFunc) (*Vase, error) {
	options := defaultOptions()
	for _, optFn := range opts {
		optFn(&options)
	}
	return &Vase{dir: dir, opts: options}, nil
}

func (v *Vase) Close() error {
	v.m.Lock()
	defer v.m.Unlock()
	if v.w != nil {
		w := v.w
		v.w = nil
		return w.Close()
	}
	return nil
}

func (v *Vase) Put(enc func(io.Writer) error) error {
	v.m.Lock()
	defer v.m.Unlock()

	if err := v.lazyInit(); err != nil {
		return err
	}

	if err := enc(v.w); err != nil {
		defer v.w.Close()
		v.w = nil
		return err
	}

	return nil
}

func (v *Vase) lazyInit() error {
	if v.w != nil {
		return nil
	}

	if err := v.opts.MkdirAll(v.dir, 0755); err != nil {
		return err
	}

	now := v.opts.Now()
	fname := fmt.Sprintf("%016x", now.UTC().UnixMicro())

	w, err := v.opts.CreateFile(path.Join(v.dir, fname))
	if err != nil {
		return err
	}

	v.w = &bufWriteCloser{bufio.NewWriter(w), w.Close}
	return nil
}

type bufWriteCloser struct {
	*bufio.Writer
	close func() error
}

func (bwc *bufWriteCloser) Close() error {
	if err := bwc.Flush(); err != nil {
		bwc.close()
		return err
	}
	return bwc.close()
}
