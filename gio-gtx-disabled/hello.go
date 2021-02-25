// SPDX-License-Identifier: Unlicense OR MIT

package main

// A simple Gio program. See https://gioui.org for more information.

import (
	"fmt"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"gioui.org/font/gofont"
)

func main() {
	go func() {
		w := app.NewWindow()
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

const (
	btn1Txt = `Click me to disable
everything for 2 seconds`
	btn2Txt = `1) Click me while enabled.
2) Hover in while disabled and click after enabled.
3) However out and in before clicking.`
)

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())
	type clickInfo struct {
		b string
		t time.Time
	}
	var (
		ops      op.Ops
		b1, b2   widget.Clickable
		disabled bool
		c        <-chan time.Time
		clicks   []clickInfo
	)
	list := layout.List{Axis: layout.Vertical}
	for {
		select {
		case <-c:
			c = nil
			disabled = false
			w.Invalidate()
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				if b1.Clicked() {
					clicks = append(clicks, clickInfo{"b1", e.Now})
					disabled = true
					c = time.After(2 * time.Second)
				}
				if b2.Clicked() {
					clicks = append(clicks, clickInfo{"b2", e.Now})
				}
				if disabled {
					gtx = gtx.Disabled()
				}
				layout.NW.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Rigid(material.Button(th, &b1, btn1Txt).Layout),
						layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
						layout.Rigid(material.Button(th, &b2, btn2Txt).Layout),
						layout.Rigid(layout.Spacer{Width: unit.Dp(5)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return list.Layout(gtx, len(clicks)+1, func(gtx layout.Context, index int) layout.Dimensions {
								msg := "I will show clicks here"
								if index != 0 {
									clk := clicks[len(clicks)-index]
									msg = fmt.Sprintf("btn %s at %s", clk.b, clk.t.Format("15:04:05"))
								}
								return material.Body1(th, msg).Layout(gtx)
							})
						}),
					)
				})
				e.Frame(gtx.Ops)
			}
		}
	}
}
