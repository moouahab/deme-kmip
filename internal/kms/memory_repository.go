package kms

import (
	"context"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	keys map[string]Key
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		keys: make(map[string]Key),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, key Key) (Key, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.keys[key.ID]; exists {
		return Key{}, ErrKeyExists
	}
	now := time.Now().UTC()
	key.CreatedAt = now
	key.UpdatedAt = now
	r.keys[key.ID] = key
	return key, nil
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (Key, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	key, exists := r.keys[id]
	if !exists {
		return Key{}, ErrKeyNotFound
	}
	return key, nil
}

func (r *MemoryRepository) Update(ctx context.Context, key Key) (Key, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.keys[key.ID]
	if !exists {
		return Key{}, ErrKeyNotFound
	}
	key.CreatedAt = existing.CreatedAt
	key.UpdatedAt = time.Now().UTC()
	r.keys[key.ID] = key
	return key, nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Key, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]Key, 0, len(r.keys))
	for _, key := range r.keys {
		keys = append(keys, key)
	}
	return keys, nil
}
