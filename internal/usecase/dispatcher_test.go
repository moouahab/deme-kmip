package usecase

import (
	"context"
	"errors"
	"testing"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

func TestDispatcherCreateKey(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	req := models.OperationRequest{
		Operation:  ttlv.OperationCreate,
		ObjectType: ttlv.ObjectTypeSymmetricKey,
	}

	resp, err := dispatcher.Dispatch(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.KeyID == "" {
		t.Fatal("expected key id, got empty string")
	}

	if resp.Status != string(kms.KeyStatusActive) {
		t.Fatalf("expected status active, got %s", resp.Status)
	}

	_, err = repo.Get(ctx, resp.KeyID)
	if err != nil {
		t.Fatalf("expected key to be stored, got %v", err)
	}

	events := auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].Operation != "create_key" {
		t.Fatalf("expected operation create_key, got %s", events[0].Operation)
	}

	if events[0].Result != "success" {
		t.Fatalf("expected result success, got %s", events[0].Result)
	}

	if events[0].KeyID != resp.KeyID {
		t.Fatalf("expected audit key id %s, got %s", resp.KeyID, events[0].KeyID)
	}
}

func TestDispatcherGetKey(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	_, err := repo.Create(ctx, kms.Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     kms.KeyStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}

	req := models.OperationRequest{
		Operation: ttlv.OperationGet,
		KeyID:     "key-123",
	}

	resp, err := dispatcher.Dispatch(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.KeyID != "key-123" {
		t.Fatalf("expected key-123, got %s", resp.KeyID)
	}

	if resp.Status != string(kms.KeyStatusActive) {
		t.Fatalf("expected status active, got %s", resp.Status)
	}

	events := auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].Operation != "get_key" {
		t.Fatalf("expected operation get_key, got %s", events[0].Operation)
	}

	if events[0].Result != "success" {
		t.Fatalf("expected result success, got %s", events[0].Result)
	}

	if events[0].KeyID != "key-123" {
		t.Fatalf("expected audit key id key-123, got %s", events[0].KeyID)
	}
}

func TestDispatcherGetKeyNotFound(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	req := models.OperationRequest{
		Operation: ttlv.OperationGet,
		KeyID:     "missing-key",
	}

	_, err := dispatcher.Dispatch(ctx, req)
	if !errors.Is(err, kms.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	events := auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].Operation != "get_key" {
		t.Fatalf("expected operation get_key, got %s", events[0].Operation)
	}

	if events[0].Result != "not_found" {
		t.Fatalf("expected result not_found, got %s", events[0].Result)
	}

	if events[0].KeyID != "missing-key" {
		t.Fatalf("expected audit key id missing-key, got %s", events[0].KeyID)
	}
}

func TestDispatcherGetKeyDestroyedIsHidden(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	_, err := repo.Create(ctx, kms.Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     kms.KeyStatusDestroyed,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}

	req := models.OperationRequest{
		Operation: ttlv.OperationGet,
		KeyID:     "key-123",
	}

	_, err = dispatcher.Dispatch(ctx, req)
	if !errors.Is(err, kms.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	events := auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].Operation != "get_key" {
		t.Fatalf("expected operation get_key, got %s", events[0].Operation)
	}

	if events[0].Result != "not_found" {
		t.Fatalf("expected result not_found, got %s", events[0].Result)
	}

	if events[0].Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected audit status destroyed, got %s", events[0].Status)
	}
}

func TestDispatcherDestroyKey(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	_, err := repo.Create(ctx, kms.Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     kms.KeyStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}

	req := models.OperationRequest{
		Operation: ttlv.OperationDestroy,
		KeyID:     "key-123",
	}

	resp, err := dispatcher.Dispatch(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resp.KeyID != "key-123" {
		t.Fatalf("expected key-123, got %s", resp.KeyID)
	}

	if resp.Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected status destroyed, got %s", resp.Status)
	}

	stored, err := repo.Get(ctx, "key-123")
	if err != nil {
		t.Fatalf("expected key to still exist internally, got %v", err)
	}

	if stored.Status != kms.KeyStatusDestroyed {
		t.Fatalf("expected stored key status destroyed, got %s", stored.Status)
	}

	events := auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].Operation != "destroy_key" {
		t.Fatalf("expected operation destroy_key, got %s", events[0].Operation)
	}

	if events[0].Result != "success" {
		t.Fatalf("expected result success, got %s", events[0].Result)
	}

	if events[0].Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected audit status destroyed, got %s", events[0].Status)
	}
}

func TestDispatcherDestroyKeyAlreadyDestroyed(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	_, err := repo.Create(ctx, kms.Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     kms.KeyStatusDestroyed,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}

	req := models.OperationRequest{
		Operation: ttlv.OperationDestroy,
		KeyID:     "key-123",
	}

	_, err = dispatcher.Dispatch(ctx, req)
	if !errors.Is(err, kms.ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	events := auditLogger.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}

	if events[0].Operation != "destroy_key" {
		t.Fatalf("expected operation destroy_key, got %s", events[0].Operation)
	}

	if events[0].Result != "not_found" {
		t.Fatalf("expected result not_found, got %s", events[0].Result)
	}

	if events[0].Status != string(kms.KeyStatusDestroyed) {
		t.Fatalf("expected audit status destroyed, got %s", events[0].Status)
	}
}

func TestDispatcherUnsupportedOperation(t *testing.T) {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	dispatcher := NewDispatcher(repo, auditLogger)

	req := models.OperationRequest{
		Operation: ttlv.Operation(999),
	}

	_, err := dispatcher.Dispatch(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	events := auditLogger.Events()
	if len(events) != 0 {
		t.Fatalf("expected 0 audit events for unsupported operation, got %d", len(events))
	}
}