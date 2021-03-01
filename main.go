package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow(app.Size(unit.Dp(800), unit.Dp(700)))

		if err := loop(w); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func logF(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	log.Print(str)
	logText.Insert(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), str))
}

func loop(w *app.Window) error {
	th := material.NewTheme(gofont.Collection())

	var permResult <-chan PermResult

	var viewEvent app.ViewEvent

	var ops op.Ops
	for {
		select {
		case result := <-permResult:
			permResult = nil
			logF("Perm result: %t %s", result.Authorized, result.Err)

			files, err := ioutil.ReadDir("/sdcard/DCIM/Camera")
			if err != nil {
				logF("read sdcard err: %s", err)
			} else {
				var names []string
				for _, f := range files {
					names = append(names, f.Name())
				}
				logF("sdcard pictures: %+v", names)
			}
			w.Invalidate()
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case app.ViewEvent:
				viewEvent = e
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)

				var reqPermClicked bool
				for btn.Clicked() {
					reqPermClicked = true
				}

				if reqPermClicked {
					permResult = RequestPermission(viewEvent)
				}

				layout.Inset{
					Bottom: e.Insets.Bottom,
					Left:   e.Insets.Left,
					Right:  e.Insets.Right,
					Top:    e.Insets.Top,
				}.Layout(gtx, func(gtx C) D {
					return drawLayout(gtx, th)
				})
				e.Frame(gtx.Ops)
			}
		}
	}
}

type (
	C = layout.Context
	D = layout.Dimensions
)

var (
	btn        = new(widget.Clickable)
	logText    = new(widget.Editor)
	layoutList = &layout.List{
		Axis: layout.Vertical,
	}
)

func drawLayout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	widgets := []layout.Widget{
		material.Button(th, btn, "Request Storage Access").Layout,
		material.Editor(th, logText, "").Layout,
	}

	return layoutList.Layout(gtx, len(widgets), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx, widgets[i])
	})
}
