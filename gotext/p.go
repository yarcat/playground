package main

import (
	"log"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// P changes printer on the fly.
type P struct {
	p *message.Printer
}

func NewP(lang string) *P {
	p := &P{}
	p.SetLang(lang)
	return p
}

func (p *P) SetLang(lang string) {
	log.Println("Parsed:", language.MustParse(lang))
	p.p = message.NewPrinter(language.Make(lang))
}

func (p *P) Printer() *message.Printer { return p.p }

func (p *P) Translate(id message.Reference, args ...interface{}) string {
	return p.p.Sprintf(id, args...)
}
