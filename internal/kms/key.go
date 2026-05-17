package kms

import (
	"time"

	"kmipDemo/internal/ttlv"
)

// KeyStatus représente l'état d'une clé dans son cycle de vie.
type KeyStatus string

const (
	KeyStatusPreActive KeyStatus = "pre_active"
	KeyStatusActive    KeyStatus = "active"
	KeyStatusRevoked   KeyStatus = "revoked"
	KeyStatusDestroyed KeyStatus = "destroyed"
)

// Key représente une clé gérée par notre KMS simulé.

// Attention : on ne stocke pas de vraie matière secrète ici.
// important : on simule la gestion de clé, mais
// on ne logge ni n'expose jamais de secret.

type Key struct {
	ID         string          `json:"id"`
	CreatedAt  time.Time       `json:"created_at"`
	ObjectType ttlv.ObjectType `json:"object_type"`
	Status     KeyStatus       `json:"status"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
