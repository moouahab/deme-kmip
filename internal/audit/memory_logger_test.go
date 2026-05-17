package audit

import (
	"context"
	"testing"
	"time"
)

func TestMemoryLoggerLogEvent(t *testing.T) {
	logger := NewMemoryLogger()
	ctx := context.Background()

	err := logger.Log(ctx, Event{
		Operation: "create_key",
		KeyID:     "key-123",
		Status:    "active",
		Result:    "success",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events := logger.Events()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	event := events[0]

	if event.Operation != "create_key" {
		t.Fatalf("expected operation create_key, got %s", event.Operation)
	}

	if event.KeyID != "key-123" {
		t.Fatalf("expected key_id key-123, got %s", event.KeyID)
	}

	if event.Status != "active" {
		t.Fatalf("expected status active, got %s", event.Status)
	}

	if event.Result != "success" {
		t.Fatalf("expected result success, got %s", event.Result)
	}

	if event.Time.IsZero() {
		t.Fatal("expected time to be set")
	}
}

func TestMemoryLoggerKeepsProvidedTime(t *testing.T) {
	logger := NewMemoryLogger()
	ctx := context.Background()

	eventTime := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)

	err := logger.Log(ctx, Event{
		Time:      eventTime,
		Operation: "get_key",
		KeyID:     "key-123",
		Result:    "success",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	events := logger.Events()

	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if !events[0].Time.Equal(eventTime) {
		t.Fatalf("expected time %v, got %v", eventTime, events[0].Time)
	}
}

func TestMemoryLoggerEventsReturnsCopy(t *testing.T) {
	logger := NewMemoryLogger()
	ctx := context.Background()

	err := logger.Log(ctx, Event{
		Operation: "destroy_key",
		KeyID:     "key-123",
		Status:    "destroyed",
		Result:    "success",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	events := logger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	events[0].Result = "modified"
	eventsAgain := logger.Events()
	if eventsAgain[0].Result != "success" {
		t.Fatalf("expected internal event to stay success, got %s", eventsAgain[0].Result)
	}
}

func TestMemoryLoggerMultipleEvents(t *testing.T) {
	logger := NewMemoryLogger()
	ctx := context.Background()

	_ = logger.Log(ctx, Event{
		Operation: "create_key",
		KeyID:     "key-1",
		Result:    "success",
	})

	_ = logger.Log(ctx, Event{
		Operation: "get_key",
		KeyID:     "key-1",
		Result:    "success",
	})

	events := logger.Events()

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	if events[0].Operation != "create_key" {
		t.Fatalf("expected first operation create_key, got %s", events[0].Operation)
	}

	if events[1].Operation != "get_key" {
		t.Fatalf("expected second operation get_key, got %s", events[1].Operation)
	}
}
