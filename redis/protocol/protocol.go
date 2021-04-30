package protocol

import (
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
