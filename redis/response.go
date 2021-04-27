package redis

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

type Response struct {
	r *bufio.Reader
}

func (r Response) Error() error {
	b, err := r.r.ReadByte()
	if err != nil {
		return err
	}
	if b == '-' {
		return r.error()
	}
	return r.r.UnreadByte()
}

func (r Response) Int() (n int, err error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return
	}
	if b != ':' {
		return 0, fmt.Errorf("%w: got=%q, want=':'", ErrFirstByte, b)
	}
	_, err = fmt.Fscanln(r.r, &n)
	return
}

func (r Response) BytesSimple() (buf []byte, err error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return
	}
	if b != '+' {
		return nil, fmt.Errorf(`%w: got=%q, want='+'`, ErrFirstByte, b)
	}
	buf, err = r.r.ReadBytes('\n')
	if l := len(buf); l > 0 && buf[l-1] == '\n' {
		buf = buf[:l-1]
	}
	return
}

func (r Response) BytesBulk(buf []byte) (n int, err error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return
	}
	if b != '$' {
		return 0, fmt.Errorf("%w: got=%q, want='$'", ErrFirstByte, b)
	}
	if _, err = fmt.Fscanln(r.r, &n); err != nil {
		return
	}
	if len(buf) > n {
		buf = buf[:n]
	}
	if _, err = io.ReadFull(r.r, buf); err != nil {
		return
	}
	if d := n - len(buf); d > 0 {
		_, err = r.r.Discard(d)
	}
	return
}

func (r Response) Consume() error { return consumeResponse(r.r) }

func (r Response) error() error {
	s := bufio.NewScanner(r.r)
	if !s.Scan() {
		return fmt.Errorf("%v: %v", ErrInternal, s.Err())
	}
	return fmt.Errorf("%w: %s", ErrResponse, s.Bytes())
}

func consumeResponse(r *bufio.Reader) error {
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch b {
	case '+', '-':
		return consumeSimpleString(r)
	case '$':
		return consumeBulkString(r)
	case '*':
		return consumeArray(r)
	case ':':
		return consumeInt(r)
	}
	return fmt.Errorf(`%w: got=%q, want one of "+-$*:"`, ErrFirstByte, b)
}

func consumeSimpleString(r *bufio.Reader) error {
	b, err := r.ReadSlice('\n')
	log.Printf("consumed %q", b)
	return err
}

func consumeBulkString(r *bufio.Reader) error {
	var n int
	if _, err := fmt.Fscanln(r, &n); err != nil {
		return err
	}
	if n <= 0 {
		return consumeSimpleString(r)
	}
	if _, err := r.Discard(n); err != nil {
		return err
	}
	return consumeSimpleString(r)
}

func consumeArray(r *bufio.Reader) error {
	return nil
}

func consumeInt(r *bufio.Reader) error {
	return nil
}
