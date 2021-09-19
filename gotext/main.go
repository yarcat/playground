package main

func main() {
	p := new(P)
	th := NewTheme(p)
	for _, lang := range []string{"en-GB", "ru-RU"} {
		p.SetLang(lang)
		c := Card{
			Messages: CardMessages{
				Title: func() string { return th.P().Sprintf("Resistivity") },
				Min:   func(x float64) string { return th.P().Sprintf("Resistivity min: %.3f", x) },
				Max:   func(x float64) string { return th.P().Sprintf("Resistivity max: %.3f", x) },
			},
		}
		c.Layout()
	}
}
