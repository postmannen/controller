package main

import (
	"context"
	"log"
	"sync"

	"gioui.org/app"
	"github.com/postmannen/controller"
)

func main() {
	eventCh := make(chan controller.Event)
	c := controller.NewController(eventCh)
	ctx := context.Background()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.Run(ctx)
		if err != nil {
			log.Printf("%v\n", err)
		}
	}()

	// Test message.
	e := controller.Event{EventType: controller.ETPrint}
	c.AddEvent(e)

	app.Main()

	wg.Wait()
}
