package protocol

import (
	"bufio"
	"io"

	"github.com/yarcat/playground/redis"
)

type (
	ArgFunc func(Writer)

	ResFunc func(Reader)

	Protocol struct {
		conn *redis.Conn
		err  error
	}
)

func Exec(p *Protocol, cmd string, args ArgFunc, res ResFunc) error {
	p.WriteString(cmd).WriteArgs(args).FinishRequest()
	return p.ReadResult(res).Err()
}

func New(conn *redis.Conn) *Protocol { return &Protocol{conn: conn} }

func (p *Protocol) Err() error { return p.err }
func (p *Protocol) OK() bool   { return p.err == nil }

func (p *Protocol) WriteString(s string) *Protocol {
	if p.err == nil {
		p.err = errOnly(p.conn.WriteString(s))
	}
	return p
}

func (p *Protocol) WriteStringEscaped(s string) *Protocol {
	if p.err == nil {
		p.err = errOnly(p.conn.WriteStringEscaped(s))
	}
	return p
}

func (p *Protocol) WriteArgs(args ArgFunc) *Protocol {
	if p.err == nil {
		args(Writer{p: p})
	}
	return p
}

func (p *Protocol) ReadResult(res ResFunc) *Protocol {
	if p.err == nil {
		res(Reader{p})
	}
	return p
}

func (p *Protocol) ConsumeMessage() *Protocol {
	if p.err == nil {
		p.err = consumeMessage(p.conn.Reader)
	}
	return p
}

func (p *Protocol) Consume(m Msg) *Protocol {
	if p.err == nil {
		p.err = consume(p.conn.Reader, m)
	}
	return p
}

func (p *Protocol) ReadInt(f func(int)) *Protocol {
	var n int
	if p.err == nil {
		p.err = parseInt(p.conn.Reader, &n)
	}
	if p.err == nil {
		f(n)
	}
	return p
}

// Read reads bulk strings into the buffer. It is expected the stream contains
// the length line and data.
func (p *Protocol) Read(buf []byte, f func(int)) *Protocol {
	if p.err != nil {
		return p
	}
	var n int
	if p.err = parseInt(p.conn.Reader, &n); p.err != nil {
		return p
	} else if n < 0 {
		f(n)
	} else if p.err = readBytes(p.conn.Reader, buf, n); p.err == nil {
		f(n)
	}
	return p
}

func (p *Protocol) FinishRequest() *Protocol {
	p.WriteString("\r\n")
	if p.err == nil {
		p.err = p.conn.Flush()
	}
	return p
}

func (p *Protocol) Msg(f func(Msg)) *Protocol {
	if p.err != nil {
		return p
	}
	var b byte
	if b, p.err = p.conn.ReadByte(); p.err != nil {
		return p
	}
	switch m := Msg(b); {
	case !m.Valid():
		p.err = ErrMsg
	default:
		f(m)
	}
	return p
}

func (p *Protocol) ReadLine(f func([]byte)) *Protocol {
	if p.err != nil {
		return p
	}
	var b []byte
	if b, p.err = p.conn.ReadSlice('\n'); p.err != nil {
		return p
	}
	f(dropCRLF(b))
	return p
}

func errOnly(n int, err error) error { return err }

func dropCRLF(b []byte) []byte {
	if len(b) < 2 {
		return b
	}
	// TODO: add checks for CRLF.
	l := len(b)
	return b[:l-2]
}

func readBytes(r *bufio.Reader, b []byte, n int) error {
	if len(b) > n {
		b = b[:n]
	}
	if _, err := io.ReadFull(r, b); err != nil {
		return err
	}
	if d := n - len(b); d > 0 {
		_, err := r.Discard(d)
		return err
	}
	_, err := r.Discard(2) // Skip CRLF
	return err
}
