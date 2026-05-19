package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

func TestHandleKMIPGetKey(t *testing.T) {
	f := newHTTPFixture()
	_, err := f.repo.Create(t.Context(), kms.Key{
		ID:         "key-123",
		ObjectType: ttlv.ObjectTypeSymmetricKey,
		Status:     kms.KeyStatusActive,
	})
	if err != nil {
		t.Fatalf("expected no error while creating key, got %v", err)
	}

	rec := serveKMIP(t, f, getKeyBody("key-123"))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}

	var resp models.OperationResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	if resp.KeyID != "key-123" {
		t.Fatalf("expected key-123, got %s", resp.KeyID)
	}
	if resp.Status != string(kms.KeyStatusActive) {
		t.Fatalf("expected active, got %s", resp.Status)
	}

	snapshot := f.metrics.Snapshot(t.Context())
	if snapshot.HTTPRequestsTotal != 1 || snapshot.GetKeyTotal != 1 || snapshot.SuccessTotal != 1 {
		t.Fatalf("unexpected metrics snapshot: %+v", snapshot)
	}
}

func TestHandleKMIPGetKeyNotFound(t *testing.T) {
	f := newHTTPFixture()
	rec := serveKMIP(t, f, getKeyBody("missing-key"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", rec.Code, rec.Body.String())
	}

	snapshot := f.metrics.Snapshot(t.Context())
	if snapshot.HTTPRequestsTotal != 1 || snapshot.GetKeyTotal != 1 || snapshot.NotFoundTotal != 1 {
		t.Fatalf("unexpected metrics snapshot: %+v", snapshot)
	}
}
