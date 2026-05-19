package usecase

import (
	"context"
	"testing"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

func TestDispatcherLocateKeysHidesDestroyed(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-active", kms.KeyStatusActive)
	createTestKey(t, f, "key-destroyed", kms.KeyStatusDestroyed)

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationLocate,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Keys) != 1 {
		t.Fatalf("expected 1 located key, got %d", len(resp.Keys))
	}
	if resp.Keys[0].ID != "key-active" {
		t.Fatalf("expected key-active, got %s", resp.Keys[0].ID)
	}
}

func TestDispatcherGetKeyAttributes(t *testing.T) {
	f := newDispatcherFixture()
	createTestKey(t, f, "key-123", kms.KeyStatusActive)

	resp, err := f.dispatcher.Dispatch(f.ctx, models.OperationRequest{
		Operation: ttlv.OperationGetAttributes,
		KeyID:     "key-123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Attributes["unique_identifier"] != "key-123" {
		t.Fatalf("expected unique_identifier key-123, got %v", resp.Attributes["unique_identifier"])
	}
	if resp.Attributes["state"] != string(kms.KeyStatusActive) {
		t.Fatalf("expected state active, got %v", resp.Attributes["state"])
	}
}

func TestDispatcherUnsupportedOperation(t *testing.T) {
	f := newDispatcherFixture()

	_, err := f.dispatcher.Dispatch(context.Background(), models.OperationRequest{
		Operation: ttlv.Operation(999),
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if events := f.auditLogger.Events(); len(events) != 0 {
		t.Fatalf("expected 0 audit events for unsupported operation, got %d", len(events))
	}
}
