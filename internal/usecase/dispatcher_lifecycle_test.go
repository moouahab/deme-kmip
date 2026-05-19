package usecase

import (
	"errors"
	"testing"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

func TestDispatcherDestroyKey(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusActive)

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationDestroy,
		KeyID:     "key-123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected status destroyed, got %s", resp.Status)
	}

	stored, err := f.repo.Get(f.ctx, "key-123")
	if err != nil {
		t.Fatalf("expected key to still exist internally, got %v", err)
	}
	if stored.Status != kms.KeyStatusDestroyed {
		t.Fatalf("expected stored key status destroyed, got %s", stored.Status)
	}

	event := requireOneAuditEvent(t, f, "destroy_key", "success")
	if event.Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected audit status destroyed, got %s", event.Status)
	}
}

func TestDispatcherDestroyKeyAlreadyDestroyed(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusDestroyed)

	_, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationDestroy,
		KeyID:     "key-123",
	})
	if !errors.Is(err, kms.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	event := requireOneAuditEvent(t, f, "destroy_key", "not_found")
	if event.Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected audit status destroyed, got %s", event.Status)
	}
}

func TestDispatcherActivateKey(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusRevoked)

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationActivate,
		KeyID:     "key-123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != string(kms.KeyStatusActive) {
		t.Fatalf("expected active, got %s", resp.Status)
	}
}

func TestDispatcherRevokeKey(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusActive)

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationRevoke,
		KeyID:     "key-123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Status != string(kms.KeyStatusRevoked) {
		t.Fatalf("expected revoked, got %s", resp.Status)
	}
}
