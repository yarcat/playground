package redis

import "bufio"

type sender struct {
	w   *bufio.Writer
	err error
}

func newSender(w *bufio.Writer) *sender { return &sender{w: w} }

func (tx *sender) String(s string) *sender {
	if tx.err == nil {
		_, tx.err = tx.w.WriteString(s)
	}
	return tx
}

var escapes = [32]string{
	`\x00`, `\x01`, `\x02`, `\x03`, `\x04`, `\x05`, `\x06`, `\x07`,
	`\x08`, `\x09`, `\n`, `\x0b`, `\r`, `\x0d`, `\x0e`, `\x0f`,
	`\x10`, `\x11`, `\x12`, `\x13`, `\x14`, `\x15`, `\x16`, `\x17`,
	`\x18`, `\x19`, `\x1a`, `\x1b`, `\x1c`, `\x1d`, `\x1e`, `\x1f`,
}

func (tx *sender) StringEscaped(s string) *sender {
	if tx.err != nil {
		return tx
	}
	for i := 0; i < len(s); {
		switch b := s[i]; {
		case b == '\\':
			tx.String(s[:i])
			tx.String(`\\`)
			s, i = s[i+1:], 0
		case b == '"':
			tx.String(s[:i])
			tx.String(`\"`)
			s, i = s[i+1:], 0
		case b < 32:
			tx.String(s[:i])
			tx.String(escapes[b])
			s, i = s[i+1:], 0
		default:
			i++
		}
	}
	if len(s) > 0 {
		tx.String(s)
	}
	return tx
}

func (tx *sender) Finish() error {
	tx.String(EOL)
	if tx.err == nil {
		tx.err = tx.w.Flush()
	}
	return tx.err
}
