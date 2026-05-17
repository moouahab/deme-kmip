package kms

import (
	"context"
	"errors"
)

var (
	ErrKeyNotFound = errors.New("kms: key not found")
	ErrKeyExists   = errors.New("kms: key already exists")
)

type Repository interface {
	Create(ctx context.Context, key Key) (Key, error)
	Get(ctx context.Context, id string) (Key, error)
	Update(ctx context.Context, key Key) (Key, error)
	List(ctx context.Context) ([]Key, error)
}
