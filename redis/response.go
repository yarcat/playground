package redis

import (
	"bufio"
	"fmt"
)

type Response struct {
	r *bufio.Reader
}

func (r Response) error() error {
	buf, err := r.r.ReadSlice('\n')
	if err != nil {
		return err
	}
	return fmt.Errorf("%w: %s", ErrStatus, dropCRLF(buf))
}

func (r Response) consumeAndErrUnexpected(b byte, want string) (err error) {
	if err = consumeResponse(r.r); err != nil {
		return err
	}
	return fmt.Errorf("%w: got=%q, want=%s", ErrFirstByte, b, want)
}

func (r Response) Int() (n int, err error) {
	switch b, err := r.r.ReadByte(); {
	case err != nil:
		return 0, err
	case b == '-':
		return 0, r.error()
	case b != ':':
		r.r.UnreadByte()
		return 0, r.consumeAndErrUnexpected(b, "':'")
	}
	return r.int()
}

func (r Response) Void() error {
	switch b, err := r.r.ReadByte(); {
	case err != nil:
		return err
	case b == '-':
		return r.error()
	case b != '+':
		r.r.UnreadByte()
		return r.consumeAndErrUnexpected(b, "'+'")
	}
	_, err := r.r.ReadSlice('\n')
	return err
}

func (r Response) int() (n int, err error) { return readInt(r.r) }

func (r Response) BytesBulk(buf []byte) (n int, err error) {
	switch b, err := r.r.ReadByte(); {
	case err != nil:
		return 0, err
	case b == '-':
		return 0, r.error()
	case b != '$':
		r.r.UnreadByte()
		return 0, r.consumeAndErrUnexpected(b, "'$'")
	}
	if n, err = r.int(); n < 0 || err != nil {
		return
	}
	return n, readBytes(r.r, buf, n)
}

func readInt(r *bufio.Reader) (n int, err error) {
	sign := 1
	if b, err := r.ReadByte(); err != nil {
		return 0, err
	} else if b == '-' {
		sign = -1
	} else {
		n = int(b - '0') // Not using unread to save one iteration.
	}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if b == '\r' {
			_, err = r.Discard(1) // Skip newline.
			return sign * n, err
		}
		n *= 10
		n += int(b - '0')
	}
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
	n, err := readInt(r)
	if err != nil {
		return err
	}
	for ; n > 0; n-- {
		if err = consumeResponse(r); err != nil {
			return err
		}
	}
	return nil
}

func consumeInt(r *bufio.Reader) error {
	_, err := readInt(r)
	return err
}
