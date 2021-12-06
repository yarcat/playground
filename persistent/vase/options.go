package vase

import (
	"io"
	"io/fs"
	"os"
	"time"
)

type (
	NowFunc        func() time.Time
	MkdirAllFunc   func(string, fs.FileMode) error
	CreateFileFunc func(string) (io.WriteCloser, error)
)

type options struct {
	Now        NowFunc
	MkdirAll   MkdirAllFunc
	CreateFile CreateFileFunc
}

type OptionFunc func(*options)

func defaultOptions() options {
	return options{
		Now:        time.Now,
		MkdirAll:   os.MkdirAll,
		CreateFile: func(name string) (io.WriteCloser, error) { return os.Create(name) },
	}
}

func WithNow(f NowFunc) OptionFunc {
	return func(o *options) { o.Now = f }
}

func WithMkdirAll(f MkdirAllFunc) OptionFunc {
	return func(o *options) { o.MkdirAll = f }
}

func WithCreateFile(f CreateFileFunc) OptionFunc {
	return func(o *options) { o.CreateFile = f }
}
