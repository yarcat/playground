package main

import (
	"errors"
	"fmt"

	"github.com/yarcat/playground/tokenizer"
	"github.com/yarcat/playground/tokenizer/examples/contrib"
)

func main() {
	for _, expr := range []string{
		"-10",
		"1 + 2*3",
		"2*3 + 4",
		"(2+2)*2",
		"-2+2*-2",
		"(2+2*2",
		"2+2)*2",
		"2*10-",
	} {
		x, err := eval(expr)
		if err == nil {
			fmt.Printf("   OK: %13s = %d\n", expr, x)
		} else {
			fmt.Printf("ERROR: %13s   %v\n", expr, err)
		}
	}
}

var (
	errUnexpectedTok = errors.New("unexpected token")
	errUnclosedParen = errors.New("unclosed parentheses")
)

type (
	token   = tokenizer.Token
	tokenFn func() token
)

func tokenGenerator(input string) tokenFn {
	t := tokenizer.New(input, contrib.DefaultNumericExpression)
	return func() tokenizer.Token {
		if !t.Next() {
			panic("must get a token")
		}
		return t.Token()
	}
}

func eval(input string) (int, error) {
	nextToken := tokenGenerator(input)
	x, tok, err := expr(nextToken)
	if err != nil {
		return 0, err
	}
	if tok.Type != tokenizer.TokenEOF {
		return 0, fmt.Errorf("%w: want EOF got %s", errUnexpectedTok, tok)
	}
	return x, nil
}

func expr(nextToken tokenFn) (int, token, error) {
	return sum(nextToken, nextToken())
}

var operatorActions = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
}

type parserFn func(tokenFn, token) (int, token, error)

func operator(nextToken tokenFn, left token, next parserFn, isOperator func(token) bool) (int, token, error) {
	res, tok, err := next(nextToken, left)
	if err != nil || tok.Type == tokenizer.TokenEOF {
		return res, tok, err
	}
	for isOperator(tok) {
		op := tok.Value
		var r int
		r, tok, err = next(nextToken, nextToken())
		if err != nil {
			return res, tok, err
		}
		res = operatorActions[op](res, r)
	}
	return res, tok, nil
}

func isSum(tok token) bool {
	return tok.Type == contrib.TokenContribOperator &&
		(tok.Value == "+" || tok.Value == "-")
}

func sum(nextToken tokenFn, left token) (int, token, error) {
	return operator(nextToken, left, mul, isSum)
}

func isMul(tok token) bool {
	return tok.Type == contrib.TokenContribOperator &&
		(tok.Value == "*" || tok.Value == "/")
}

func mul(nextToken tokenFn, left token) (int, token, error) {
	return operator(nextToken, left, una, isMul)
}

func isUnary(tok token) bool {
	return tok.Type == contrib.TokenContribOperator &&
		(tok.Value == "+" || tok.Value == "-")
}

func una(nextToken tokenFn, tok token) (int, token, error) {
	if !isUnary(tok) {
		return fac(nextToken, tok)
	}
	op := tok.Value
	x, tok, err := fac(nextToken, nextToken())
	if op == "-" {
		x = -x
	}
	return x, tok, err
}

func fac(nextToken tokenFn, tok token) (int, token, error) {
	switch tok.Type {
	case contrib.TokenContribOpenParen:
		x, tok, err := expr(nextToken)
		if err != nil {
			return 0, tok, err
		}
		if tok.Type != contrib.TokenContribClosingParen {
			return 0, tok, errUnclosedParen
		}
		return x, nextToken(), nil
	case contrib.TokenContribNumber:
		var x int
		fmt.Sscan(tok.Value, &x)
		return x, nextToken(), nil
	}
	return 0, tok, fmt.Errorf("%w: %s", errUnexpectedTok, tok)
}
