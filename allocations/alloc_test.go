package allocations

import (
	"bufio"
	"testing"
)

type fakeStream struct{}

func (fs fakeStream) Read(b []byte) (int, error)  { return len(b), nil }
func (fs fakeStream) Write(b []byte) (int, error) { return len(b), nil }

var (
	c     = newClient(fakeStream{})
	value = []byte("value")
	key   = "key"
)

func BenchmarkSet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.Set("key", "value")
	}
}

func BenchmarkExecFuncBufio(b *testing.B) {
	b.Run("Factory", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.ExecFuncBufio("set", String(key), Bytes(value))
		}
	})

	b.Run("Lambda", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			c.ExecFuncBufio("set",
				func(w *bufio.Writer) error { return writeString(w, key) },
				func(w *bufio.Writer) error { return writeBytes(w, value) },
			)
		}
	})

	b.Run("Functors", func(b *testing.B) {
		b.Run("Ptr", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				c.ExecFuncBufio("set",
					StringWriter{&key}.WriteArg,
					BytesWriter{&value}.WriteArg,
				)
			}
		})
		b.Run("Val", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				c.ExecFuncBufio("set", String2(key), Bytes2(value))
			}
		})
		b.Run("Hlp", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Set(c).
					KV(String2(key), Bytes2(value)).
					KV(String2("key"), Bytes2(value)).
					Exec()
			}
		})
	})
}
func BenchmarkExecI(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.ExecInterface("set",
			StringWriter{&key},
			BytesWriter{&value},
		)
	}
}

func BenchmarkExecFuncSenderPtr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.ExecFuncSenderPtr("set", func(args *sender) {
			args.String("key").Bytes(value)
		})
	}
}

func BenchmarkExecFuncSenderVal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.ExecFuncSenderVal("set", func(args sender) sender {
			return *args.String("key").Bytes(value)
		})
	}
}

func BenchmarkExecFuncSenderPrealloc(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c.ExecFuncSenderPrealloc("set", func(args *sender) {
			args.String("key").Bytes(value)
		})
	}
}
