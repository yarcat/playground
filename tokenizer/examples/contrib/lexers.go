package contrib

import (
	"strings"
	"unicode"

	"github.com/yarcat/playground/tokenizer"
)

// RuneMatcherFn matches the rune.
type RuneMatcherFn func(rune) bool

// NumericExpression tokenizes numeric expressions e.g. "(2+2)*2".
type NumericExpression struct {
	IsDigit    RuneMatcherFn
	IsOperator RuneMatcherFn
	IsSpace    RuneMatcherFn
}

var DefaultNumericExpression = NumericExpression{
	IsDigit:    unicode.IsDigit,
	IsOperator: func(r rune) bool { return strings.ContainsRune("+-*/", r) },
	IsSpace:    unicode.IsSpace,
}.Scan

// Scan parses and emits numbers, operators or parentheses.
func (ne NumericExpression) Scan(s *tokenizer.Lexer) tokenizer.StateFn {
	switch r, eof := s.Next(); {
	case eof:
		s.Emit(tokenizer.TokenEOF)
		return nil
	case r == '(':
		s.Emit(TokenContribOpenParen)
	case r == ')':
		s.Emit(TokenContribClosingParen)
	case unicode.IsDigit(r):
		s.AcceptRunFn(unicode.IsDigit)
		s.Emit(TokenContribNumber)
	case strings.ContainsRune("-+*/", r):
		s.Emit(TokenContribOperator)
	case unicode.IsSpace(r):
		s.Ignore()
	default:
		return s.Errorf("unexpected %q", r)
	}
	return ne.Scan

}
