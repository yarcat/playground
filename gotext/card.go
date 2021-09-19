package main

import (
	"fmt"

	"golang.org/x/text/message"
)

type Card struct {
	Messages CardMessages
}

func (c Card) Layout() {
	fmt.Printf(`
Title: %s
Min:   %s
Max:   %s
`,
		c.Messages.Title(),
		c.Messages.Min(123),
		c.Messages.Max(321),
	)
}

type CardMessages struct {
	Title func() string
	Min   func(float64) string
	Max   func(float64) string
}

type CardMessagesBuilder struct {
	Title string
	Min   string
	Max   string
}

func (cmb CardMessagesBuilder) Build(p *message.Printer) CardMessages {
	return CardMessages{
		Title: func() string { return p.Sprintf(cmb.Title) },
		Min:   func(f float64) string { return p.Sprintf(cmb.Min, f) },
		Max:   func(f float64) string { return p.Sprintf(cmb.Max, f) },
	}
}
