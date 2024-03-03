# Event Emitter

`event_emitter` is a Go package that provides a simple event emitter implementation.

## Installation

You can install the package using `go get`:

```bash
go get github.com/anboo/event-emitter
```

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"github.com/anboo/event-emitter"
)

func main() {
	// Create a new event emitter
	emitter := event_emitter.NewEventEmitter()

	// Subscribe to an event
	event_emitter.Subscribe(func(event string) {
		fmt.Println("Received event:", event)
	})

	// Emit an event
	emitter.Emit("Hello, world!")
}
```

### Using Custom Event Types

```go
package main

import (
	"fmt"
	"github.com/anboo/event-emitter"
)

// CustomEvent is a custom event type
type CustomEvent struct {
	Message string
}

func main() {
	// Create a new event emitter
	emitter := event_emitter.NewEventEmitter()

	// Subscribe to an event with a custom type
	event_emitter.Subscribe(func(event CustomEvent) {
		fmt.Println("Received custom event:", event.Message)
	})

	// Emit an event with a custom type
	customEvent := CustomEvent{"Hello, custom world!"}
	emitter.Emit(customEvent)
}
```

### Unsubscribing from Events

```go
package main

import (
	"fmt"
	"github.com/anboo/event-emitter"
)

func main() {
	// Create a new event emitter
	emitter := event_emitter.NewEventEmitter()

	// Subscribe to an event
	subscription := func(event string) {
		fmt.Println("Received event:", event)
	}
	event_emitter.Subscribe(subscription)

	// Emit an event
	emitter.Emit("Hello, world!")

	// Unsubscribe from the event
	event_emitter.Unsubscribe(subscription)

	// Emit the event again
	emitter.Emit("Hello again, world!") // This should not trigger any subscriber
}
```

## License
