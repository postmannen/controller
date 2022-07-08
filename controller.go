package controller

import (
	"context"
	"fmt"
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
						c.errCh <- fmt.Errorf("got etDone")
					}()
				}
			case <-ctx.Done():
				c.errCh <- fmt.Errorf("info: got ctx.Done")
			}
		}
	}()

	// Split this out in an ErrorKernel
	for {
		err := <-c.errCh
		return err
	}
}

func (c *controller) AddEvent(event Event) {
	c.eventCh <- event
}
