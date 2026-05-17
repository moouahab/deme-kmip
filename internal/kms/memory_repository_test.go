package kms

import (
	"context"
	"errors"
	"testing"

	"kmipDemo/internal/ttlv"
)

func TestMemoryRepositoryCreateAndGet(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     KeyStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got, err := repo.Get(ctx, "key-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != created.ID {
		t.Fatalf("expected id %s, got %s", created.ID, got.ID)
	}

	if got.Status != KeyStatusActive {
		t.Fatalf("expected status active, got %s", got.Status)
	}
}

func TestMemoryRepositoryCreateDuplicate(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	key := Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     KeyStatusActive,
	}

	_, err := repo.Create(ctx, key)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.Create(ctx, key)
	if !errors.Is(err, ErrKeyExists) {
		t.Fatalf("expected ErrKeyExists, got %v", err)
	}
}

func TestMemoryRepositoryGetNotFound(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.Get(ctx, "missing-key")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestMemoryRepositoryUpdate(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.Create(ctx, Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     KeyStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := repo.Update(ctx, Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     KeyStatusDestroyed,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Status != KeyStatusDestroyed {
		t.Fatalf("expected status destroyed, got %s", updated.Status)
	}
}

func TestMemoryRepositoryList(t *testing.T) {
	repo := NewMemoryRepository()
	ctx := context.Background()

	_, _ = repo.Create(ctx, Key{
		ID:         "key-1",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     KeyStatusActive,
	})

	_, _ = repo.Create(ctx, Key{
		ID:         "key-2",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     KeyStatusActive,
	})

	keys, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}
