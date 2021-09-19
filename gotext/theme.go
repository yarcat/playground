package main

import "golang.org/x/text/message"

type Theme struct {
	p *P
}

func NewTheme(p *P) *Theme { return &Theme{p} }

func (t *Theme) P() *message.Printer { return t.p.Printer() }
func (t *Theme) Translate(id message.Reference, args ...interface{}) string {
	return t.p.Translate(id, args)
}
