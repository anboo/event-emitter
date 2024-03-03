# Event Emitter

`event_emitter` is a Go package that provides a simple event emitter implementation.

## Installation

You can install the package using `go get`:

```bash
go get github.com/anboo/event-emitter
```

## Usage

### Basic Usage global emitter

```go
package main

import (
	"fmt"

	"github.com/anboo/event-emitter"
)

func main() {
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
```

### Local emitter
```go
package main

import (
	"fmt"

	"github.com/anboo/event-emitter"
)

func main() {
	emitter := event_emitter.NewEventEmitter()

	// Подписываемся на событие типа ExampleEvent.
	event_emitter.Subscribe(func(event ExampleEvent) {
		fmt.Println("Received event:", event.Message)
	}, emitter)

	event_emitter.Emit(ExampleEvent{"Hello, world!"}, emitter)
}

// ExampleEvent - пример структуры события.
type ExampleEvent struct {
	Message string
}
```

## License
