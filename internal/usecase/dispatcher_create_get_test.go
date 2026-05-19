package usecase

import (
	"errors"
	"testing"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

func TestDispatcherCreateKey(t *testing.T) {
	f := newDispatcherFixture()

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation:  ttlv.OperationCreate,
		ObjectType: ttlv.ObjectTypeSymmetricKey,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.KeyID == "" {
		t.Fatal("expected key id, got empty string")
	}
	if resp.Status != string(kms.KeyStatusActive) {
		t.Fatalf("expected status active, got %s", resp.Status)
	}
	if _, err := f.repo.Get(f.ctx, resp.KeyID); err != nil {
		t.Fatalf("expected key to be stored, got %v", err)
	}

	event := requireOneAuditEvent(t, f, "create_key", "success")
	if event.KeyID != resp.KeyID {
		t.Fatalf("expected audit key id %s, got %s", resp.KeyID, event.KeyID)
	}
}

func TestDispatcherGetKey(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusActive)

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationGet,
		KeyID:     "key-123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.KeyID != "key-123" {
		t.Fatalf("expected key-123, got %s", resp.KeyID)
	}
	if resp.Status != string(kms.KeyStatusActive) {
		t.Fatalf("expected status active, got %s", resp.Status)
	}

	event := requireOneAuditEvent(t, f, "get_key", "success")
	if event.KeyID != "key-123" {
		t.Fatalf("expected audit key id key-123, got %s", event.KeyID)
	}
}

func TestDispatcherGetKeyNotFound(t *testing.T) {
	f := newDispatcherFixture()

	_, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationGet,
		KeyID:     "missing-key",
	})
	if !errors.Is(err, kms.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	event := requireOneAuditEvent(t, f, "get_key", "not_found")
	if event.KeyID != "missing-key" {
		t.Fatalf("expected audit key id missing-key, got %s", event.KeyID)
	}
}

func TestDispatcherGetKeyDestroyedIsHidden(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusDestroyed)

	_, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationGet,
		KeyID:     "key-123",
	})
	if !errors.Is(err, kms.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	event := requireOneAuditEvent(t, f, "get_key", "not_found")
	if event.Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected audit status destroyed, got %s", event.Status)
	}
}
