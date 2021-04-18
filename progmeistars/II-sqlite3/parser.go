package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	Input string
	tag   string
}

func firstTag(s string) (tag string, pos int) {
	pos = -1
	for _, t := range []string{"[artist]", "[title]", "[genre]", "##", "&&&"} {
		if i := strings.Index(s, t); i >= 0 && (i < pos || pos < 0) {
			tag, pos = t, i
		}
	}
	return
}

func (p *Parser) Next() bool {
	tag, pos := firstTag(p.Input)
	if pos < 0 {
		return false
	}
	p.Input = p.Input[pos+len(tag):]
	p.tag = tag
	return true
}

var tagAliases = map[string]string{
	"&&&": "tracks",
	"##":  "year",
}

func (p *Parser) Tag() string {
	if tag, ok := tagAliases[p.tag]; ok {
		return tag
	}
	return strings.Trim(p.tag, "[]")
}

func (p *Parser) Text() string {
	_, pos := firstTag(p.Input)
	if pos < 0 {
		pos = len(p.Input)
	}
	return strings.TrimSpace(p.Input[:pos])
}

func (p *Parser) Int() (n int) {
	fmt.Sscan(p.Text(), &n)
	return
}
