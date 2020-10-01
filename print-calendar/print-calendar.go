package main

import (
	"fmt"
	"strings"
	"time"
)

func month(y int, m time.Month) (c [][]string) {
	week := make([]string, 7)
	wrap := func() {
		c = append(c, week)
		week = make([]string, 7)
	}
	addDay := func(t time.Time) (wrapped bool) {
		d := (6 + int(t.Weekday())) % 7 // Our week starts on Monday.
		week[d] = fmt.Sprintf("%d", t.Day())
		if d == len(week)-1 { // Last day of the week filled -- wrap the week.
			wrap()
			wrapped = true
		}
		return

	}
	var wrapped bool
	t0 := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	for t := t0; t0.Month() == t.Month(); t = t.AddDate(0, 0, 1 /* days */) {
		wrapped = addDay(t)
	}
	if !wrapped {
		wrap()
	}
	return
}

func main() {
	for m := 1; m <= 12; m++ {
		cal := month(1980, time.Month(m))
		fmt.Printf("%15s | Mon Tue Wed Thu Fri Sat Sun\n", time.Month(m))
		for _, week := range cal {
			fmt.Printf("%s | ", strings.Repeat(" ", 15))
			for _, day := range week {
				fmt.Printf("%3s ", day)
			}
			fmt.Println()
		}
	}
}
