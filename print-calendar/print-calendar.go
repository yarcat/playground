package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func monthCalendar(y int, m time.Month) (c [][]string) {
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

var date = flag.String("date", "today",
	"YYYYMM, YYYY, MM format date or 'today' to print the calendar for")

func dateFromFlag() (year int, month time.Month, monthOk bool, err error) {
	var d time.Time
	if *date == "today" {
		d = time.Now()
		year, month, monthOk = d.Year(), d.Month(), true
	} else if len(*date) == 6 { // YYYYMM
		if d, err = time.Parse("200601", *date); err != nil {
			return
		}
		year, month, monthOk = d.Year(), d.Month(), true
	} else if len(*date) == 4 { // YYYY
		if d, err = time.Parse("2006", *date); err != nil {
			return
		}
		year = d.Year()
	} else if len(*date) == 2 { // MM
		if d, err = time.Parse("01", *date); err != nil {
			return
		}
		year, month, monthOk = time.Now().Year(), d.Month(), true
	} else {
		err = fmt.Errorf("expected to be 'today', YYYYMM, YYYY or MM, but got %q", *date)
	}
	return
}

func printMonthCalendar(year int, month time.Month) {
	cal := monthCalendar(year, month)
	fmt.Printf("%15s | Mon Tue Wed Thu Fri Sat Sun\n", month)
	for _, week := range cal {
		fmt.Printf("%s | ", strings.Repeat(" ", 15))
		for _, day := range week {
			fmt.Printf("%3s ", day)
		}
		fmt.Println()
	}
}

func main() {
	flag.Parse()
	year, month, singleMonth, err := dateFromFlag()
	if err != nil {
		fmt.Printf("Error parsing -date: %v\n", err)
		os.Exit(1)
	}
	if singleMonth {
		printMonthCalendar(year, month)
		return
	}
	for m := 1; m <= 12; m++ {
		printMonthCalendar(year, time.Month(m))
	}
}
