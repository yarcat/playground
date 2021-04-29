package redis

import (
	"bufio"
	"bytes"
	"io"
)

// Stream represents a single connection with a Redis service. The stream is
// not go-routine safe.
type Stream struct {
	r      *bufio.Reader
	sender sender
}

// NewStream initializes a stream. There should be only one stream per
// connection to a server.
func NewStream(stream io.ReadWriter) *Stream {
	return &Stream{
		r:      bufio.NewReader(stream),
		sender: sender{w: bufio.NewWriter(stream)},
	}
}

// Sender provides methods for sending data through a stream.
type Sender struct{ s *sender }

// ByteRaw sends the byte as is.
func (send Sender) ByteRaw(b byte) {
	// TODO(yarcat): move this to the sender struct.
	if send.s.err == nil {
		send.s.err = send.s.w.WriteByte(b)
	}
}

// String sends the string provided in a binary-safe mode.
func (send Sender) String(s string) { send.s.StringEscaped(s) }

// Bytes sends the bytes provided in a binary-safe mode.
// func (send Sender) Bytes(b []byte) { send.s.Bytes(b) }

// ExecuteArgFunc is an argument type passed to Execute. The function is
// responsible for serializing an argument using the argument sender provided.
type ExecuteArgFunc func(Sender)

// String should wrap string arguments.
func String(s string) ExecuteArgFunc { return sendString{s}.Send }

type sendString struct{ s string }

func (ss sendString) Send(send Sender) {
	send.ByteRaw('"')
	send.String(ss.s)
	send.ByteRaw('"')
}

// Execute is a request builder that runs the request and returns the result object.
// Execution result must be consumed.
func Execute(stream *Stream, cmd string, args ...ExecuteArgFunc) ExecuteResult {
	sender := Sender{&stream.sender}
	sender.String(cmd)
	for _, a := range args {
		sender.ByteRaw(' ')
		a(sender)
	}
	return ExecuteResult{r: stream.r, err: sender.s.Finish()}
}

// ArrayItem returns ExecuteResult without sending requests. The function is
// meant for consuming array items.
func ArrayItem(stream *Stream) ExecuteResult { return ExecuteResult{r: stream.r} }

// Status contains errors received from Redis responses.
type Status struct{ b []byte }

// Bytes returns the status/error bytes read from the response.
func (er Status) Bytes() []byte { return er.b }

// OK returns true if there is no error.
func (er Status) OK() bool {
	return len(er.b) == 0 || bytes.Equal(er.Bytes(), []byte("OK"))
}

// ExecuteResult is returned by Execute and must be used to consume the
// response from the Redis stream.
type ExecuteResult struct {
	r   *bufio.Reader
	err error
}

var msgUnexpectedResponse = []byte("received unexpected response")

// Void represents no result. It silently consumes OK status.
func (er ExecuteResult) Void() (status Status, err error) {
	if er.err != nil {
		return status, er.err
	}
	switch b, err := er.r.ReadByte(); {
	case err != nil:
		return status, err
	case b == '-' || b == '+':
		buf, err := er.r.ReadSlice('\n')
		return Status{dropCRLF(buf)}, err
	}
	if err = er.r.UnreadByte(); err != nil {
		return status, err
	}
	return Status{msgUnexpectedResponse}, consumeResponse(er.r)
}

// Int expects a numeric result.
func (er ExecuteResult) Int(n *int) (status Status, err error) {
	if er.err != nil {
		return status, er.err
	}
	switch b, err := er.r.ReadByte(); {
	case err != nil:
		return status, err
	case b == ':':
		*n, err = readInt(er.r)
		return status, err
	}
	if err = er.r.UnreadByte(); err != nil {
		return status, err
	}
	return Status{msgUnexpectedResponse}, consumeResponse(er.r)
}

// Bytes expects a bulk string.
func (er ExecuteResult) Bytes(data []byte, n *int, found *bool) (status Status, err error) {
	if er.err != nil {
		return status, err
	}
	switch b, err := er.r.ReadByte(); {
	case err != nil:
		return status, err
	case b == '$':
		switch *n, err = readInt(er.r); {
		case err != nil:
			return status, err
		case *n < 0:
			*n, *found = 0, false
			return status, nil
		}
		*found = true
		return status, readBytes(er.r, data, *n)
	}
	if err = er.r.UnreadByte(); err != nil {
		return status, err
	}
	return Status{msgUnexpectedResponse}, consumeResponse(er.r)
}

// Array expects an array.
func (er ExecuteResult) Array(n *int) (status Status, err error) {
	if er.err != nil {
		return status, er.err
	}
	switch b, err := er.r.ReadByte(); {
	case err != nil:
		return status, err
	case b != '*':
		if err = er.r.UnreadByte(); err != nil {
			return status, err
		}
		return Status{msgUnexpectedResponse}, consumeResponse(er.r)
	}
	*n, err = readInt(er.r)
	return status, err
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

func dropCRLF(b []byte) []byte {
	if l := len(b); l >= 2 {
		return b[:l-2]
	}
	return b
}
