package main

import (
	"fmt"

	"github.com/anboo/event-emitter"
)

func main() {
	event_emitter.Emitter.Debug()

	event_emitter.Subscribe(func(event ExampleEvent) {
		fmt.Println("Received event:", event.Message)
	})

	event_emitter.Subscribe(func(event AnotherEvent) {
		fmt.Println("Received another event:", event.Data)
	})

	event1 := ExampleEvent{"Hello, world!"}
	event_emitter.Emit(event1)

	event2 := AnotherEvent{42}
	event_emitter.Emit(event2)

	event_emitter.Unsubscribe(AnotherEvent{})

	event_emitter.Emit(event2)
}

type ExampleEvent struct {
	Message string
}

type AnotherEvent struct {
	Data int
}
