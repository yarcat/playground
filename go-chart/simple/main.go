package main

import (
	"image"
	"io"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func timeseries(dt time.Duration, f func(time.Time) float64) (x []time.Time, y []float64) {
	now := time.Now()
	for t := now.Add(-24 * time.Hour); now.After(t); t = t.Add(dt) {
		x = append(x, t)
		y = append(y, f(t))
	}
	return x, y
}

func fixed(v float64) func(time.Time) float64 {
	return func(time.Time) float64 { return v }
}

func random() func(time.Time) float64 {
	var (
		dv    float64
		steps int
	)
	v := 50.0
	var lastT time.Time
	return func(t time.Time) float64 {
		if lastT.IsZero() {
			lastT = t
			return v
		}
		dtSec := t.Sub(lastT).Seconds()
		lastT = t
		if steps == 0 {
			targetV := v - 5 + rand.Float64()*10
			fillInSteps := math.Abs(targetV-v) * (3600 / dtSec) / 100 // 100% in one hour
			steps = int(fillInSteps*0.5 + rand.Float64()*fillInSteps) // +/- 100%
			if steps == 0 {
				steps = 1
			}
			dv = (targetV - v) / float64(steps)
		}
		v += dv
		if v < 0 {
			v = 0
		} else if v > 100 {
			v = 100
		}
		steps--

		return v
	}
}

func makeLine(x []time.Time, c drawing.Color, v float64) chart.TimeSeries {
	return chart.TimeSeries{
		XValues: []time.Time{x[0], x[len(x)-1]},
		YValues: []float64{v, v},
		Style: chart.Style{
			StrokeColor: c,
			StrokeWidth: 2,
		},
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const (
		dt                = 60 * time.Second
		widthPx, heightPx = 200, 200
	)
	x, y := timeseries(dt, fixed(50))
	//x, y := timeseries(dt, random())

	log.Printf(`Plotting sample graph:
	number of points  : %v
	interval:         : %v
	width, height (px): %v %v
	`, len(x), dt, widthPx, heightPx)

	main := chart.TimeSeries{
		XValues: x,
		YValues: y,
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
		},
	}
	mainOutline := chart.TimeSeries{
		XValues: x,
		YValues: y,
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
			StrokeWidth: 2,
		},
	}

	graph := chart.Chart{
		Width:  widthPx,
		Height: heightPx,
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeHourValueFormatter,
			GridMajorStyle: chart.Style{
				StrokeWidth: 1.5,
				StrokeColor: chart.ColorAlternateGray.WithAlpha(80),
			},
			GridMinorStyle: chart.Style{
				StrokeWidth: 1,
				StrokeColor: chart.ColorAlternateGray.WithAlpha(40),
			},
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{Min: 0, Max: 100},
			GridMajorStyle: chart.Style{
				StrokeWidth: 1,
				StrokeColor: chart.ColorAlternateGray.WithAlpha(80),
			},
			GridMinorStyle: chart.Style{
				StrokeWidth: 0.5,
				StrokeColor: chart.ColorAlternateGray.WithAlpha(40),
			},
		},
		Series: []chart.Series{
			main,
			mainOutline,
			makeLine(main.XValues, chart.ColorAlternateGreen, 30),
			makeLine(main.XValues, chart.ColorAlternateGreen, 70),
			makeLine(main.XValues, chart.ColorAlternateYellow, 15),
			makeLine(main.XValues, chart.ColorAlternateYellow, 85),
			makeLine(main.XValues, chart.ColorRed, 5),
			makeLine(main.XValues, chart.ColorRed, 95),
		},
	}

	start := time.Now()
	graph.Render(chart.PNG, collector{})
	// var b bytes.Buffer
	// graph.Render(chart.PNG, &b)
	// os.WriteFile("out.png", b.Bytes(), 0644)
	log.Println(time.Since(start))
}

type collector struct{ io.Writer }

func (c collector) SetRGBA(img *image.RGBA) {
	sz := img.Bounds().Max
	log.Println("collecting image width, height (px): ", sz.X, sz.Y)
}
