package main

import (
	"github.com/yarcat/playground/tokenizer"
	"github.com/yarcat/playground/tokenizer/examples/contrib"
)

var zeroToken tokenizer.Token

type scanner struct {
	t          *tokenizer.Tokenizer
	prev, last tokenizer.Token
}

func newScanner(input string) *scanner {
	return &scanner{
		t: tokenizer.New(input, contrib.DefaultNumericExpression),
	}
}

func (s *scanner) MustToken() tokenizer.Token {
	if s.prev != zeroToken {
		s.last, s.prev = s.prev, zeroToken
		return s.last
	}
	if !s.t.Next() {
		panic("must get a token")
	}
	s.last = s.t.Token()
	return s.last
}

func (s *scanner) Backout() { s.prev = s.last }
