package allocations

import (
	"bufio"
	"io"
	"testing"
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

type fakeStream struct{}

func (fs fakeStream) Read(b []byte) (int, error)  { return len(b), nil }
func (fs fakeStream) Write(b []byte) (int, error) { return len(b), nil }

func BenchmarkSet(b *testing.B) {
	c := newClient(fakeStream{})

	b.Run("Set", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Set("key", "value")
		}
	})
}

type ExecArgFunc func(*bufio.Writer) error

func writeString(w *bufio.Writer, s string) error {
	_, err := w.WriteString(s)
	return err
}

func String(s string) ExecArgFunc {
	return func(w *bufio.Writer) error { return writeString(w, s) }
}

func writeBytes(w *bufio.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}

func Bytes(b []byte) ExecArgFunc {
	return func(w *bufio.Writer) error { return writeBytes(w, b) }
}

func (c client) Exec(cmd string, exec ...ExecArgFunc) error {
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

func BenchmarkExec(b *testing.B) {
	c := newClient(fakeStream{})
	value := []byte("value")
	key := "key"

	b.Run("Factory", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Exec("set", String(key), Bytes(value))
		}
	})

	b.Run("Lambda", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Exec("set",
				func(w *bufio.Writer) error { return writeString(w, key) },
				func(w *bufio.Writer) error { return writeBytes(w, value) },
			)
		}
	})

}

func (c client) Exec2(cmd string, args func(*sender)) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	sender := newSender(c.w)
	args(sender)
	return sender.Send()
}

func BenchmarkExec2(b *testing.B) {
	c := newClient(fakeStream{})
	value := []byte("value")

	b.Run("Exec2", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Exec2("set", func(args *sender) {
				args.String("key").Bytes(value)
			})
		}
	})
}

func (c client) Exec3(cmd string, args func(sender) sender) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	sender := args(*newSender(c.w))
	return sender.Send()
}

func BenchmarkExec3(b *testing.B) {
	c := newClient(fakeStream{})
	value := []byte("value")

	b.Run("Exec3", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Exec3("set", func(args sender) sender {
				return *args.String("key").Bytes(value)
			})
		}
	})
}

func (c client) Exec4(cmd string, args func(*sender)) error {
	_, err := c.w.WriteString(cmd)
	if err != nil {
		return err
	}
	args(c.s)
	return c.s.Send()
}

func BenchmarkExec4(b *testing.B) {
	c := newClient(fakeStream{})
	value := []byte("value")

	b.Run("Exec4", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.Exec4("set", func(args *sender) {
				args.String("key").Bytes(value)
			})
		}
	})
}
