package redis

import (
	"bufio"
	"io"
)

type Conn struct{ bufio.ReadWriter }

func NewConn(stream io.ReadWriter) *Conn {
	return &Conn{ReadWriter: *bufio.NewReadWriter(
		bufio.NewReader(stream),
		bufio.NewWriter(stream),
	)}
}

func (conn *Conn) WriteStringEscaped(s string) (int, error) {
	return newSender(conn.Writer).StringEscaped(s).Result()
}
