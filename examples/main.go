package main

import (
	"fmt"

	"github.com/anboo/event-emitter"
)

func main() {
	// Подписываемся на событие типа ExampleEvent.
	event_emitter.Subscribe(func(event ExampleEvent) {
		fmt.Println("Received event:", event.Message)
	})

	// Подписываемся на событие типа AnotherEvent.
	event_emitter.Subscribe(func(event AnotherEvent) {
		fmt.Println("Received another event:", event.Data)
	})

	// Создаем и отправляем события разных типов.
	event1 := ExampleEvent{"Hello, world!"}
	event_emitter.Emit(event1)

	event2 := AnotherEvent{42}
	event_emitter.Emit(event2)

	// Отписываемся от событий типа AnotherEvent.
	event_emitter.Unsubscribe(AnotherEvent{})

	// Отправляем еще одно событие типа AnotherEvent.
	event_emitter.Emit(event2)
}

// ExampleEvent - пример структуры события.
type ExampleEvent struct {
	Message string
}

// AnotherEvent - другой пример структуры события.
type AnotherEvent struct {
	Data int
}
