package httpapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/usecase/models"
)

func TestHandleKMIPCreateKey(t *testing.T) {
	f := newHTTPFixture()
	rec := serveKMIP(t, f, createKeyBody())

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected application/json, got %s", rec.Header().Get("Content-Type"))
	}

	var resp models.OperationResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	if resp.KeyID == "" {
		t.Fatal("expected key_id in response, got empty string")
	}

	stored, err := f.repo.Get(t.Context(), resp.KeyID)
	if err != nil {
		t.Fatalf("expected key to be stored in repository, got %v", err)
	}
	if stored.Status != kms.KeyStatusActive {
		t.Fatalf("expected stored key status active, got %s", stored.Status)
	}

	snapshot := f.metrics.Snapshot(t.Context())
	if snapshot.HTTPRequestsTotal != 1 || snapshot.CreateKeyTotal != 1 || snapshot.SuccessTotal != 1 {
		t.Fatalf("unexpected metrics snapshot: %+v", snapshot)
	}
}
