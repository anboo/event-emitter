package event_emitter

import (
	"log/slog"
	"os"
	"reflect"
	"sync"
)

// Emitter is the default global instance of EventEmitter.
var Emitter = NewEventEmitter("global")

// EventEmitter represents an event emitter.
type EventEmitter struct {
	subscribers map[reflect.Type][]func(event interface{}) // Map of event types to subscriber functions
	mutex       sync.RWMutex

	loggerLevel *slog.LevelVar
	logger      *slog.Logger
}

func (e *EventEmitter) Debug() {
	e.loggerLevel.Set(slog.LevelDebug)
}

// NewEventEmitter creates a new instance of EventEmitter.
func NewEventEmitter(name string, loggerOptionalParam ...*slog.Logger) *EventEmitter {
	var (
		loggerLevel *slog.LevelVar
		logger      *slog.Logger
	)

	if len(loggerOptionalParam) == 0 {
		loggerLevel = &slog.LevelVar{}
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: loggerLevel,
		})).With("service", "event-emitter").With("name", name)
	} else {
		logger = loggerOptionalParam[0]
	}

	return &EventEmitter{
		subscribers: make(map[reflect.Type][]func(event interface{})),
		mutex:       sync.RWMutex{},

		loggerLevel: loggerLevel,
		logger:      logger,
	}
}

// Subscribe subscribes to events of type T with the given subscriber function.
func Subscribe[T any](subscriber func(event T), optionalEmitter ...*EventEmitter) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	emitter := getEmitter(optionalEmitter...)

	if emitter.logger != nil {
		emitter.logger.Debug("subscribe", "event", eventType.Name())
	}

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
	if emitter.logger != nil {
		emitter.logger.Debug("unsubscribe", "event", eventType.Name())
	}

	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	emitter.subscribers[eventType] = make([]func(event interface{}), 0)
}

// getEmitter returns the appropriate emitter instance.
func getEmitter(optionalEmitterParam ...*EventEmitter) *EventEmitter {
	if len(optionalEmitterParam) == 0 {
		return Emitter
	}
	return optionalEmitterParam[0]
}

// Emit sends the event to all subscribers of the appropriate type.
func Emit(event interface{}, optionalEmitterParam ...*EventEmitter) {
	eventType := reflect.TypeOf(event)

	emitter := getEmitter(optionalEmitterParam...)

	emitter.mutex.RLock()
	defer emitter.mutex.RUnlock()

	if emitter.logger != nil {
		emitter.logger.Debug("emit", "event", eventType.Name(), "subscribers", len(emitter.subscribers[eventType]))
	}

	if subscribers, ok := emitter.subscribers[eventType]; ok {
		for _, subscriber := range subscribers {
			subscriber(event)
		}
	}
}
