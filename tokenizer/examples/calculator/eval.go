package main

import (
	"errors"
	"fmt"

	"github.com/yarcat/playground/tokenizer"
	"github.com/yarcat/playground/tokenizer/examples/contrib"
)

var errUnexpectedTok = errors.New("unexpected token")

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

func expr(s *scanner) (int, error) { return sum(s) }

func sum(s *scanner) (int, error) {
	return operator(s, mul, isSum)
}

func mul(s *scanner) (int, error) {
	return operator(s, unary, isMul)
}

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
	tok := s.MustToken()
	switch tok.Type {
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
	}
	return 0, unexpTok(tok, "NUMBER or OPEN_PAREN")
}

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
	return backout(s, x)
}

func isSum(tok tokenizer.Token) bool {
	return tok.Type == contrib.TokenContribOperator &&
		(tok.Value == "+" || tok.Value == "-")
}

func isMul(tok tokenizer.Token) bool {
	return tok.Type == contrib.TokenContribOperator &&
		(tok.Value == "*" || tok.Value == "/")
}

func isUnaryOp(tok tokenizer.Token) bool {
	return tok.Type == contrib.TokenContribOperator &&
		(tok.Value == "+" || tok.Value == "-")
}

func backout(s *scanner, x int) (int, error) {
	s.Backout()
	return x, nil
}
