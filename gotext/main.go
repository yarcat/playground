package main

func main() {
	p := new(P)
	th := NewTheme(p)
	for _, lang := range []string{"en-GB", "ru-RU"} {
		p.SetLang(lang)
		c := Card{
			Messages: CardMessages{
				Title: func() string { return th.Translate("Resistivity") },
				Min:   func(x float64) string { return th.Translate("Resistivity min: %.3f", x) },
				Max:   func(x float64) string { return th.Translate("Resist. max: %.3f", x) },
			},
		}
		c.Layout()
	}
}
