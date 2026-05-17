package audit

import (
	"context"
	"sync"
	"time"
)

type MemoryLogger struct {
	mu     sync.RWMutex
	events []Event
}

func NewMemoryLogger() *MemoryLogger {
	return &MemoryLogger{
		events: make([]Event, 0),
	}
}

func (l *MemoryLogger) Log(ctx context.Context, event Event) error {
	_ = ctx

	l.mu.Lock()
	defer l.mu.Unlock()
	if event.Time.IsZero() {
		event.Time = time.Now().UTC()
	}
	l.events = append(l.events, event)
	return nil
}

func (l *MemoryLogger) Events() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()

	events := make([]Event, len(l.events))
	copy(events, l.events)

	return events
}
