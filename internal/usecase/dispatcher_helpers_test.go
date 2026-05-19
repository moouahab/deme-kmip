package usecase

import (
	"context"
	"testing"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
)

type dispatcherFixture struct {
	ctx         context.Context
	repo        *kms.MemoryRepository
	auditLogger *audit.MemoryLogger
	dispatcher  *Dispatcher
}

func newDispatcherFixture() dispatcherFixture {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()

	return dispatcherFixture{
		ctx:         context.Background(),
		repo:        repo,
		auditLogger: auditLogger,
		dispatcher:  NewDispatcher(repo, auditLogger),
	}
}

func createTestKey(t *testing.T, f dispatcherFixture, id string, status kms.KeyStatus) {
	t.Helper()

	_, err := f.repo.Create(f.ctx, kms.Key{
		ID:         id,
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     status,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}
}

func requireOneAuditEvent(t *testing.T, f dispatcherFixture, operation string, result string) audit.Event {
	t.Helper()

	events := f.auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	if events[0].Operation != operation {
		t.Fatalf("expected operation %s, got %s", operation, events[0].Operation)
	}
	if events[0].Result != result {
		t.Fatalf("expected result %s, got %s", result, events[0].Result)
	}

	return events[0]
}
