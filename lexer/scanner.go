package tokenizer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Scanner implements rune scanning operations.

// Every successful (no EOF) Scanner.Next() call advances current position in
// the input string, allowing to emit current buffer as a token.
type Scanner struct {
	input string
	emit  func(Token)

	start, pos int
	width      int // Last read rune size in bytes.
}

// NewScanner returns new scanner initialized with the string.
func NewScannerString(s string, emit func(Token)) *Scanner {
	return &Scanner{
		input: s,
		emit:  emit,
	}
}

// Emit produces a token of the given type and current value. The scanner resets
// start to current position in the input string.
func (s *Scanner) Emit(t TokenType) {
	tok := Token{Type: t, Value: s.input[s.start:s.pos]}
	s.start = s.pos
	s.emit(tok)
}

// Errorf emits an error oken and returns nil state function.
func (s *Scanner) Errorf(format string, args ...interface{}) StateFn {
	msg := fmt.Sprintf(format, args...)
	tok := Token{
		Type:  TokenError,
		Value: fmt.Sprintf("error at position %d: %s", s.start, msg),
	}
	s.start = s.pos
	s.emit(tok)
	return nil
}

// Next returns next rune read.
func (s *Scanner) Next() (r rune, eof bool) {
	if s.pos >= len(s.input) {
		return r, true /*eof=*/
	}
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	return r, false /*eof=*/
}

// Backout unreads the last rune. Only a single rune can be unread.
func (s *Scanner) Backout() {
	s.pos -= s.width
	s.width = 0
}

// Ignore makes the scanner to forget current buffer.
func (s *Scanner) Ignore() { s.start = s.pos }

// Accept consumes a single rune from the matching set. The function returns
// true if the rune matches, or false otherwise. If the rune wasn't matched,
// then Backout is used to return the rune back to the stream.
func (s *Scanner) Accept(set string) bool {
	r, eof := s.Next()
	if eof {
		return false
	}
	if strings.ContainsRune(set, r) {
		return true
	}
	s.Backout()
	return false
}

// AcceptRun consumes a sequence of runes by calling Accept in a loop.
func (s *Scanner) AcceptRun(set string) {
	for s.Accept(set) {
	}
}

// AcceptFn consumes a single rune. The function returns true if the filter
// returns true. If the rune wasn't filtered, then Backout method is used to
// return the rune back to the stream.
func (s *Scanner) AcceptFn(f func(rune) bool) bool {
	r, eof := s.Next()
	if eof {
		return false
	}
	if f(r) {
		return true
	}
	s.Backout()
	return false
}

// AcceptRunFn consumes a sequence of runes by calling AcceptFn in a loop.
func (s *Scanner) AcceptRunFn(f func(rune) bool) {
	for s.AcceptFn(f) {
	}
}
