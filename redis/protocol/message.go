package protocol

import (
	"bufio"
	"errors"
)

type Msg byte

const (
	MsgSimpleStr Msg = '+'
	MsgErr       Msg = '-'
	MsgInt       Msg = ':'
	MsgBulkStr   Msg = '$'
	MsgArray     Msg = '*'
)

var (
	ErrMsg = errors.New("unknown message type")
	ErrEOL = errors.New("unexpected character waiting for EOL")
)

var validMsg = []Msg{MsgSimpleStr, MsgErr, MsgInt, MsgBulkStr, MsgArray}

func (m Msg) Valid() bool {
	for _, vm := range validMsg {
		if m == vm {
			return true
		}
	}
	return false
}

func consumeMessage(r *bufio.Reader) error {
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	return consume(r, Msg(b))
}

func consume(r *bufio.Reader, m Msg) error {
	switch m {
	case MsgErr, MsgSimpleStr, MsgInt:
		return consumeSimpleStr(r)
	case MsgBulkStr:
		return consumeBulkStr(r)
	case MsgArray:
		return consumeArray(r)
	}
	return ErrMsg

}

func consumeSimpleStr(r *bufio.Reader) error {
	_, err := r.ReadSlice('\n')
	return err
}

func consumeBulkStr(r *bufio.Reader) error {
	var n int
	if err := parseInt(r, &n); err != nil {
		return err
	}
	n += 2 // Include CRLF.
	_, err := r.Discard(n)
	return err
}

func consumeArray(r *bufio.Reader) error {
	var n int
	if err := parseInt(r, &n); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		if err := consumeMessage(r); err != nil {
			return err
		}
	}
	return nil
}

func consumeCRLF(r *bufio.Reader) error {
	for _, want := range []byte{'\r', '\n'} {
		switch got, err := r.ReadByte(); {
		case err != nil:
			return err
		case got != want:
			return ErrEOL
		}
	}
	return nil
}

func errorOrConsumeCRLF(err error, r *bufio.Reader) error {
	if err != nil {
		return err
	}
	return consumeCRLF(r)
}

func parseInt(r *bufio.Reader, n *int) error {
	sign := 1
	switch b, err := r.ReadByte(); {
	case err != nil:
		return err
	case b == '-':
		sign = -1
	default:
		if err = r.UnreadByte(); err != nil {
			return err
		}
	}
	for {
		switch b, err := r.ReadByte(); {
		case err != nil:
			return err
		case !('0' <= b && b <= '9'): // Isn't a digit.
			*n *= sign
			return errorOrConsumeCRLF(r.UnreadByte(), r)
		default:
			*n *= 10
			*n += int(b - '0')
		}
	}
}
