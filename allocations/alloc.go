package allocations

import (
	"bufio"
	"io"
)

type client struct {
	r *bufio.Reader
	w *bufio.Writer
	s *sender
}

func newClient(stream io.ReadWriter) client {
	w := bufio.NewWriter(stream)
	return client{
		r: bufio.NewReader(stream),
		w: w,
		s: newSender(w),
	}
}

func (c client) Set(key, value string) (err error) {
	return newSender(c.w).String("set").Arg(key).Arg(value).Send()
}

type sender struct {
	w   *bufio.Writer
	err error
}

func newSender(w *bufio.Writer) *sender { return &sender{w: w} }

func (send *sender) String(s string) *sender {
	if send.err == nil {
		_, send.err = send.w.WriteString(s)
	}
	return send
}

func (send *sender) Arg(s string) *sender { return send.String(s) }

func (send *sender) Bytes(b []byte) *sender {
	if send.err == nil {
		_, send.err = send.w.Write(b)
	}
	return send
}

func (send *sender) Send() error {
	send.String("\r\n")
	if send.err != nil {
		send.err = send.w.Flush()
	}
	return send.err
}

type ExecArgFunc func(*bufio.Writer) error

func writeString(w *bufio.Writer, s string) error {
	_, err := w.WriteString(s)
	return err
}

func String(s string) ExecArgFunc { return func(w *bufio.Writer) error { return writeString(w, s) } }

type StringWriter struct{ str *string }

func (sw StringWriter) WriteArg(w *bufio.Writer) error { return writeString(w, *sw.str) }

func writeBytes(w *bufio.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}

func Bytes(b []byte) ExecArgFunc { return func(w *bufio.Writer) error { return writeBytes(w, b) } }

type BytesWriter struct{ bytes *[]byte }

func (bw BytesWriter) WriteArg(w *bufio.Writer) error { return writeBytes(w, *bw.bytes) }

func (c client) ExecFuncBufio(cmd string, exec ...ExecArgFunc) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	for _, f := range exec {
		if err = f(c.w); err != nil {
			return err
		}
	}
	_, err = c.w.WriteString("\r\n")
	if err == nil {
		err = c.w.Flush()
	}
	return err
}

// StringWriter2 uses value unline StringWriter, which uses a pointer.
type StringWriter2 struct{ s string }

func (sw StringWriter2) WriteArg(w *bufio.Writer) error {
	_, err := w.WriteString(sw.s)
	return err
}

// BytesWriter2 uses value unline BytesWriter, which uses a pointer.
type BytesWriter2 struct{ b []byte }

func (bw BytesWriter2) WriteArg(w *bufio.Writer) error {
	_, err := w.Write(bw.b)
	return err
}

func String2(s string) ExecArgFunc { return StringWriter2{s}.WriteArg }
func Bytes2(b []byte) ExecArgFunc  { return BytesWriter2{b}.WriteArg }

type setter struct{ s *sender }

func Set(c client) setter { return setter{c.s.String("set")} }
func (s setter) KV(key, value ExecArgFunc) setter {
	s.s.String(" ")
	key(s.s.w)
	s.s.String(" ")
	value(s.s.w)
	return s
}
func (s setter) Exec() error { return s.s.Send() }

type ArgWriter interface {
	WriteArg(b *bufio.Writer) error
}

func (c client) ExecInterface(cmd string, args ...ArgWriter) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	for _, arg := range args {
		if err = arg.WriteArg(c.w); err != nil {
			return err
		}
	}
	_, err = c.w.WriteString("\r\n")
	if err == nil {
		err = c.w.Flush()
	}
	return err
}

func (c client) ExecFuncSenderPtr(cmd string, args func(*sender)) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	sender := newSender(c.w)
	args(sender)
	return sender.Send()
}

func (c client) ExecFuncSenderVal(cmd string, args func(sender) sender) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	sender := args(*newSender(c.w))
	return sender.Send()
}

func (c client) ExecFuncSenderPrealloc(cmd string, args func(*sender)) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	args(c.s)
	return c.s.Send()
}
