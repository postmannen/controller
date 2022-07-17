package controller

import (
	"context"
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Event struct {
	EventType eventType
}

type eventType int

const (
	ETPrint eventType = iota
	ETExit
	ETDone
)

type controller struct {
	errCh   chan error
	eventCh chan Event
}

func NewController(eventCh chan Event) *controller {
	c := controller{
		errCh:   make(chan error, 1),
		eventCh: eventCh,
	}

	return &c
}

func (c *controller) Run(ctx context.Context) error {
	go func() {
		for {
			select {
			case e := <-c.eventCh:
				switch e.EventType {
				case ETPrint:
					go func() {
						fmt.Printf("info: got event: %v\n", e)
						c.eventCh <- Event{EventType: ETDone}
					}()
				case ETDone:
					go func() {
						fmt.Printf("info: got event: %v\n", e)
						c.errCh <- fmt.Errorf("info: got etDone")
					}()
				}
			case <-ctx.Done():
				c.errCh <- fmt.Errorf("info: got ctx.Done")
			}
		}
	}()

	startUI()

	// Split this out in an ErrorKernel
	for {
		err := <-c.errCh
		log.Printf("just logging error: %v\n", err)
	}

}

func (c *controller) AddEvent(event Event) {
	c.eventCh <- event
}

func startUI() {
	// Start UI.
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Controller"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)

		// ops are the operations from the UI
		var ops op.Ops

		// startButton is a clickable widget
		var startButton widget.Clickable
		var stopButton widget.Clickable

		// th defnes the material design style
		th := material.NewTheme(gofont.Collection())

		// listen for events in the window.
		for e := range w.Events() {

			// detect what type of event
			switch e := e.(type) {

			// this is sent when the application should re-render.
			case system.FrameEvent:
				{
					gtx := layout.NewContext(&ops, e)
					layout.Flex{
						Axis:    layout.Vertical,
						Spacing: layout.SpaceStart,
					}.Layout(gtx,
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &startButton, "Start")
								return btn.Layout(gtx)
							},
						),
						layout.Rigid(
							func(gtx layout.Context) layout.Dimensions {
								btn := material.Button(th, &stopButton, "stop")
								return btn.Layout(gtx)
							},
						),
						layout.Rigid(layout.Spacer{Height: (25)}.Layout),
					)

					e.Frame(gtx.Ops)
				}
			}
		}
		os.Exit(0)
	}()
}
