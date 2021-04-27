package redis

import (
	"bufio"
	"fmt"
	"io"
)

type Response struct {
	r *bufio.Reader
}

func (r Response) Error() error {
	b, err := r.r.ReadByte()
	if err != nil {
		return err
	}
	if b != '-' {
		r.r.UnreadByte()
		return nil
	}
	// It's ok to create a scanner here, since error cases
	buf, err := r.r.ReadSlice('\n')
	if err != nil {
		return fmt.Errorf("%v: %v", ErrInternal, err)
	}
	if l := len(buf); l > 0 && buf[l-1] == '\n' {
		buf = buf[:l-1]
	}
	return fmt.Errorf("%w: %s", ErrResponse, buf)
}

func (r Response) Int() (n int, err error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return
	}
	if b != ':' {
		return 0, fmt.Errorf("%w: got=%q, want=':'", ErrFirstByte, b)
	}
	return r.int()
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

func (r Response) int() (n int, err error) {
	sign := 1
	if b, err := r.r.ReadByte(); err != nil {
		return 0, err
	} else if b == '-' {
		sign = -1
	} else {
		r.r.UnreadByte()
	}
	for {
		b, err := r.r.ReadByte()
		if err != nil {
			return 0, err
		}
		if b == '\r' {
			_, err = r.r.Discard(1) // Skip newline.
			return sign * n, err
		}
		n *= 10
		n += int(b - '0')
	}
}

func (r Response) BytesBulk(buf []byte) (n int, err error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if b != '$' {
		return 0, fmt.Errorf("%w: got=%q, want='$'", ErrFirstByte, b)
	}
	if n, err = r.int(); n < 0 || err != nil {
		return
	}
	if len(buf) > n {
		buf = buf[:n]
	}
	if _, err = io.ReadFull(r.r, buf); err != nil {
		return 0, err
	}
	if d := n - len(buf); d > 0 {
		n = len(buf)
		if _, err = r.r.Discard(d); err != nil {
			return n, err
		}
	}
	_, err = r.r.Discard(2) // Skip EOL.
	return n, err
}

func (r Response) ErrorOrConsume() error {
	if err := r.Error(); err != nil {
		return err
	}
	return consumeResponse(r.r)
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
	_, err := r.ReadSlice('\n')
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
