package protocol

type Writer struct{ p *Protocol }

func (w Writer) WriteString(s string) Writer {
	w.p.WriteString(s)
	return w
}

func WriteStrings(values ...string) ArgFunc { return stringsWriter(values).Send }

type stringsWriter []string

func (ss stringsWriter) Send(w Writer) {
	for _, s := range ss {
		w.WriteString(` "`).p.WriteStringEscaped(s).WriteString(`"`)
	}
}
