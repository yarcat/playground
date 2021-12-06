package vase

import (
	"bufio"
)

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

type bufReadCloser struct {
	*bufio.Reader
	close func() error
}

func (bwc *bufReadCloser) Close() error { return bwc.close() }
