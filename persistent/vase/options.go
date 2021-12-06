package vase

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type (
	NowFunc        func() time.Time
	MkdirAllFunc   func(string, fs.FileMode) error
	CreateFileFunc func(string) (io.WriteCloser, error)
	OpenFileFunc   func(string) (io.ReadCloser, error)
	GlobFunc       func(string) ([]string, error)
)

type options struct {
	Now        NowFunc
	MkdirAll   MkdirAllFunc
	Glob       GlobFunc
	CreateFile CreateFileFunc
	OpenFile   OpenFileFunc
}

type OptionFunc func(*options)

func defaultOptions() options {
	return options{
		Now:        time.Now,
		MkdirAll:   os.MkdirAll,
		Glob:       filepath.Glob,
		CreateFile: func(name string) (io.WriteCloser, error) { return os.Create(name) },
		OpenFile:   func(name string) (io.ReadCloser, error) { return os.Open(name) },
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

func WithGlob(f GlobFunc) OptionFunc {
	return func(o *options) { o.Glob = f }
}
