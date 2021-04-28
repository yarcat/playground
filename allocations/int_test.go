package allocations

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func benchmark(b *testing.B, f func()) {
	for n := 0; n < b.N; n++ {
		f()
	}
}

var (
	byts = []byte("123\r\n")
	b1   = bytes.NewBuffer(byts)
	bb   = bufio.NewReader(b1)
)

func initStreams() {
	b1.Reset()
	b1.Write(byts)
	bb.Reset(b1)
}

func impl1() {
	initStreams()
	var n int
	fmt.Fscanln(b1, &n)
}

func impl2() {
	initStreams()
	b, _ := bb.ReadSlice('\r')
	var n int
	fmt.Sscanln(string(b), &n)
}

func impl3() {
	initStreams()
	b, _ := bb.ReadSlice('\r')
	b = b[:len(b)-1] // Cut CR.

	strconv.ParseInt(
		// Soooo dirty!!! But we really don't need any copies here.
		*(*string)(unsafe.Pointer(&reflect.StringHeader{
			Data: (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data,
			Len:  len(b),
		})),
		10, 64)
}

func impl4() {
	initStreams()
	var n int
	for {
		switch b, _ := bb.ReadByte(); b {
		case '\r':
			bb.Discard(1) // Skip LF.
			return
		default:
			n *= 10
			n += int(b - '0')
		}

	}
}

func BenchmarkImpl1(b *testing.B) { benchmark(b, impl1) }
func BenchmarkImpl2(b *testing.B) { benchmark(b, impl2) }
func BenchmarkImpl3(b *testing.B) { benchmark(b, impl3) }
func BenchmarkImpl4(b *testing.B) { benchmark(b, impl4) }
