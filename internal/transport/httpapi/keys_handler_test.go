package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
)

func TestHandleKeys(t *testing.T) {
	ctx := context.Background()
	repo := kms.NewMemoryRepository()
	_, err := repo.Create(ctx, kms.Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     kms.KeyStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}

	handler := HandleKeys(repo)
	req := httptest.NewRequest(http.MethodGet, "/keys", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var keys []kms.Key
	if err := json.NewDecoder(rec.Body).Decode(&keys); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
	if keys[0].ID != "key-123" {
		t.Fatalf("expected key-123, got %s", keys[0].ID)
	}
}
