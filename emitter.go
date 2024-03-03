package event_emitter

import (
	"reflect"
	"sync"
)

// globalEmitter is the default global instance of EventEmitter.
var globalEmitter = NewEventEmitter()

// EventEmitter represents an event emitter.
type EventEmitter struct {
	subscribers map[reflect.Type][]func(event interface{}) // Map of event types to subscriber functions
	mutex       sync.RWMutex                               // Mutex for thread safety
}

// NewEventEmitter creates a new instance of EventEmitter.
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		subscribers: make(map[reflect.Type][]func(event interface{})),
	}
}

// Subscribe subscribes to events of type T with the given subscriber function.
func Subscribe[T any](subscriber func(event T), optionalEmitter ...*EventEmitter) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	emitter := getEmitter(optionalEmitter...)

	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	if _, ok := emitter.subscribers[eventType]; !ok {
		emitter.subscribers[eventType] = make([]func(event interface{}), 0)
	}
	emitter.subscribers[eventType] = append(emitter.subscribers[eventType], func(event interface{}) {
		subscriber(event.(T))
	})
}

// Unsubscribe unsubscribes all subscribers from events of type T.
func Unsubscribe[T any](event T, optionalEmitter ...*EventEmitter) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	emitter := getEmitter(optionalEmitter...)

	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	emitter.subscribers[eventType] = make([]func(event interface{}), 0)
}

// getEmitter returns the appropriate emitter instance.
func getEmitter(optionalEmitterParam ...*EventEmitter) *EventEmitter {
	if len(optionalEmitterParam) == 0 {
		return globalEmitter
	}
	return optionalEmitterParam[0]
}

// Emit sends the event to all subscribers of the appropriate type.
func Emit(event interface{}, optionalEmitterParam ...*EventEmitter) {
	eventType := reflect.TypeOf(event)

	emitter := getEmitter(optionalEmitterParam...)
	emitter.mutex.RLock()
	defer emitter.mutex.RUnlock()

	if subscribers, ok := emitter.subscribers[eventType]; ok {
		for _, subscriber := range subscribers {
			subscriber(event)
		}
	}
}
