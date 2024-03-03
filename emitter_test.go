package event_emitter

import (
	"log/slog"
	"os"
	"testing"
)

func BenchmarkEmit(b *testing.B) {
	type ExampleEvent struct {
		Message string
	}

	//or slog.New(slog.NewTextHandler(io.Discard, nil))
	emitter := NewEventEmitter("benchmark", nil)

	// Subscribe to the event
	Subscribe(func(event ExampleEvent) {}, emitter)
	Subscribe(func(event ExampleEvent) {}, emitter)
	Subscribe(func(event ExampleEvent) {}, emitter)

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		Emit(ExampleEvent{}, emitter)
	}
}

func TestGlobalEmitter(t *testing.T) {
	emitter := Emitter

	if emitter == nil {
		t.Errorf("NewEventEmitter() returned nil")
	}

	if len(emitter.subscribers) != 0 {
		t.Errorf("New event emitter should have no subscribers")
	}

	if emitter.logger == nil {
		t.Errorf("New event emitter should have a logger")
	}
}

func TestNewEventEmitter(t *testing.T) {
	emitter := NewEventEmitter("test")

	if emitter == nil {
		t.Errorf("NewEventEmitter() returned nil")
	}

	if len(emitter.subscribers) != 0 {
		t.Errorf("New event emitter should have no subscribers")
	}

	if emitter.logger == nil {
		t.Errorf("New event emitter should have a logger")
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: &slog.LevelVar{},
	})).With("service", "event-emitter")

	emitter = NewEventEmitter("test", logger)
	if emitter.logger != logger {
		t.Errorf("New event emitter should have specific logger")
	}
}

func TestSubscribeUnsubscribeAndEmit(t *testing.T) {
	type ExampleEvent struct {
		Message string
	}

	type AnotherEvent struct {
		Data int
	}

	var counterExampleEventReached, counterAnotherEventReached int

	Subscribe(func(event ExampleEvent) {
		counterExampleEventReached++
	})

	Subscribe(func(event AnotherEvent) {
		counterAnotherEventReached++
	})

	event1 := ExampleEvent{"Hello, world!"}
	event2 := AnotherEvent{42}
	Emit(event1)
	Emit(event2)
	Emit(event2)

	if counterExampleEventReached != 1 {
		t.Fatalf("expected counterExampleEventReached is 1 got %d", counterExampleEventReached)
	}

	if counterAnotherEventReached != 2 {
		t.Fatalf("expected counterAnotherEventReached is 1 got %d", counterAnotherEventReached)
	}

	Unsubscribe(AnotherEvent{})
	Emit(event1)
	Emit(event2)

	if counterExampleEventReached != 2 {
		t.Fatalf("expected counterExampleEventReached is 2 got %d", counterExampleEventReached)
	}

	if counterAnotherEventReached != 2 {
		t.Fatalf("expected counterAnotherEventReached is 2 after unsubscribe got %d", counterAnotherEventReached)
	}
}

func TestSubscribeAndEmitStandardCases(t *testing.T) {
	type testCase struct {
		emitterSubscribe      *EventEmitter
		emitterEmit           *EventEmitter
		name                  string
		event                 string
		expectedEventReceived bool
	}

	testCases := []testCase{
		{
			name:                  "Event received correctly for global emitter",
			emitterSubscribe:      Emitter,
			emitterEmit:           Emitter,
			event:                 "test message",
			expectedEventReceived: true,
		},
		{
			name:                  "Event not received with different emitters",
			emitterSubscribe:      Emitter,
			emitterEmit:           NewEventEmitter("different"),
			event:                 "test message",
			expectedEventReceived: false,
		},
		{
			name:                  "Classic case with nil logger different emitters",
			emitterSubscribe:      Emitter,
			emitterEmit:           NewEventEmitter("different", nil),
			event:                 "test message",
			expectedEventReceived: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eventReceived := false

			// Subscribe to the event
			Subscribe(func(event string) {
				if event != "test message" {
					t.Errorf("Received unexpected event message. Expected: %s, Got: %s", "test message", event)
				}
				eventReceived = true
			}, tc.emitterSubscribe)

			// Emit the event
			Emit(tc.event, tc.emitterEmit)

			if eventReceived != tc.expectedEventReceived {
				t.Errorf("expected recevied status %v got %v", tc.expectedEventReceived, eventReceived)
			}
		})
	}
}
