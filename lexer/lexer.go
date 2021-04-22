package tokenizer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Lexer implements rune scanning operations.
//
// Every successful (no EOF) Lexer.Next() call advances current position in
// the input string, allowing to emit current buffer as a token.
type Lexer struct {
	input string
	emit  func(Token)

	start, pos int
	width      int // Last read rune size in bytes.
}

// NewScanner returns new scanner initialized with the string.
func NewScannerString(s string, emit func(Token)) *Lexer {
	return &Lexer{
		input: s,
		emit:  emit,
	}
}

// Emit produces a token of the given type and current value. The scanner resets
// start to current position in the input string.
func (l *Lexer) Emit(t TokenType) {
	tok := Token{Type: t, Value: l.input[l.start:l.pos]}
	l.start = l.pos
	l.emit(tok)
}

// Errorf emits an error oken and returns nil state function.
func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	msg := fmt.Sprintf(format, args...)
	tok := Token{
		Type:  TokenError,
		Value: fmt.Sprintf("error at position %d: %s", l.start, msg),
	}
	l.start = l.pos
	l.emit(tok)
	return nil
}

// Next returns next rune read.
func (l *Lexer) Next() (r rune, eof bool) {
	if l.pos >= len(l.input) {
		return r, true /*eof=*/
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r, false /*eof=*/
}

// Backout unreads the last rune. Only a single rune can be unread.
func (l *Lexer) Backout() {
	l.pos -= l.width
	l.width = 0
}

// Ignore makes the scanner to forget current buffer.
func (l *Lexer) Ignore() { l.start = l.pos }

// Accept consumes a single rune from the matching set. The function returns
// true if the rune matches, or false otherwise. If the rune wasn't matched,
// then Backout is used to return the rune back to the stream.
func (l *Lexer) Accept(set string) bool {
	r, eof := l.Next()
	if eof {
		return false
	}
	if strings.ContainsRune(set, r) {
		return true
	}
	l.Backout()
	return false
}

// AcceptRun consumes a sequence of runes by calling Accept in a loop.
func (l *Lexer) AcceptRun(set string) {
	for l.Accept(set) {
	}
}

// AcceptFn consumes a single rune. The function returns true if the filter
// returns true. If the rune wasn't filtered, then Backout method is used to
// return the rune back to the stream.
func (l *Lexer) AcceptFn(f func(rune) bool) bool {
	r, eof := l.Next()
	if eof {
		return false
	}
	if f(r) {
		return true
	}
	l.Backout()
	return false
}

// AcceptRunFn consumes a sequence of runes by calling AcceptFn in a loop.
func (l *Lexer) AcceptRunFn(f func(rune) bool) {
	for l.AcceptFn(f) {
	}
}
