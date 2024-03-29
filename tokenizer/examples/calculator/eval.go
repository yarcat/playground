package main

import (
	"errors"
	"fmt"

	"github.com/yarcat/playground/tokenizer"
	"github.com/yarcat/playground/tokenizer/examples/contrib"
)

var errUnexpectedTok = errors.New("unexpected token")

var operatorActions = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

type (
	actionFn  func(*scanner) (int, error)
	opMatchFn func(tokenizer.Token) bool
)

func unexpTok(got tokenizer.Token, want string) error {
	return fmt.Errorf("%w: got = %v, want = %s", errUnexpectedTok, got, want)
}

func eval(s *scanner) (int, error) {
	x, err := expr(s)
	if err != nil {
		return 0, err
	}
	if tok := s.MustToken(); tok.Type != tokenizer.TokenEOF {
		return 0, unexpTok(tok, "EOF")
	}
	return x, nil
}

func isSum(tok tokenizer.Token) bool     { return expectedOp(tok, "+", "-") }
func isMul(tok tokenizer.Token) bool     { return expectedOp(tok, "*", "/") }
func isUnaryOp(tok tokenizer.Token) bool { return expectedOp(tok, "+", "-") }

func expr(s *scanner) (int, error) { return sum(s) }
func sum(s *scanner) (int, error)  { return operator(s, mul, isSum) }
func mul(s *scanner) (int, error)  { return operator(s, unary, isMul) }

func unary(s *scanner) (int, error) {
	tok := s.MustToken()
	if !isUnaryOp(tok) {
		s.Backout()
		return factor(s)
	}
	x, err := factor(s)
	if tok.Value == "-" {
		x = -x
	}
	return x, err
}

func factor(s *scanner) (x int, err error) {
	switch tok := s.MustToken(); tok.Type {
	case contrib.TokenContribOpenParen:
		if x, err = expr(s); err != nil {
			return 0, err
		}
		if tok = s.MustToken(); tok.Type != contrib.TokenContribClosingParen {
			return 0, unexpTok(tok, "CLOSING_PAREN")
		}
		return x, nil
	case contrib.TokenContribNumber:
		fmt.Sscan(tok.Value, &x)
		return x, nil
	default:
		return 0, unexpTok(tok, "NUMBER or OPEN_PAREN")
	}
}

func operator(s *scanner, action actionFn, opMatch opMatchFn) (x int, err error) {
	if x, err = action(s); err != nil {
		return 0, err
	}
	for tok := s.MustToken(); opMatch(tok); tok = s.MustToken() {
		var y int
		if y, err = action(s); err != nil {
			return 0, err
		}
		x = operatorActions[tok.Value](x, y)
	}
	s.Backout()
	return x, nil
}

func expectedOp(tok tokenizer.Token, set ...string) bool {
	if tok.Type != contrib.TokenContribOperator {
		return false
	}
	for _, s := range set {
		if tok.Value == s {
			return true
		}
	}
	return false
}
