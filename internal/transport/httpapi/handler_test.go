package httpapi

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase"
	"kmipDemo/internal/usecase/models"
)

func TestHandleKMIPCreateKey(t *testing.T) {
	operationValue := make([]byte, 4)
	binary.BigEndian.PutUint32(operationValue, uint32(ttlv.OperationCreate))

	body := []byte{
		0x42, 0x00, 0x5C,
		0x05,
		0x00, 0x00, 0x00, 0x04,
	}
	body = append(body, operationValue...)

	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()

	dispatcher := usecase.NewDispatcher(repo, auditLogger)
	handler := HandleKMIP(dispatcher, metricsCollector)

	req := httptest.NewRequest(http.MethodPost, "/kmip", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

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

	stored, err := repo.Get(req.Context(), resp.KeyID)
	if err != nil {
		t.Fatalf("expected key to be stored in repository, got %v", err)
	}

	if stored.Status != kms.KeyStatusActive {
		t.Fatalf("expected stored key status active, got %s", stored.Status)
	}

	snapshot := metricsCollector.Snapshot(req.Context())
	if snapshot.HTTPRequestsTotal != 1 {
		t.Fatalf("expected 1 http request, got %d", snapshot.HTTPRequestsTotal)
	}
	if snapshot.CreateKeyTotal != 1 {
		t.Fatalf("expected 1 create key, got %d", snapshot.CreateKeyTotal)
	}
	if snapshot.SuccessTotal != 1 {
		t.Fatalf("expected 1 success, got %d", snapshot.SuccessTotal)
	}
}

func TestHandleKMIPGetKey(t *testing.T) {
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

	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()

	dispatcher := usecase.NewDispatcher(repo, auditLogger)
	handler := HandleKMIP(dispatcher, metricsCollector)

	operationValue := make([]byte, 4)
	binary.BigEndian.PutUint32(operationValue, uint32(ttlv.OperationGet))

	body := []byte{
		0x42, 0x00, 0x5C,
		0x05,
		0x00, 0x00, 0x00, 0x04,
	}
	body = append(body, operationValue...)

	keyIDValue := []byte("key-123")
	keyIDBlock := []byte{
		0x42, 0x00, 0x94,
		0x07,
		0x00, 0x00, 0x00, byte(len(keyIDValue)),
	}
	keyIDBlock = append(keyIDBlock, keyIDValue...)
	body = append(body, keyIDBlock...)

	req := httptest.NewRequest(http.MethodPost, "/kmip", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

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

	snapshot := metricsCollector.Snapshot(req.Context())
	if snapshot.HTTPRequestsTotal != 1 {
		t.Fatalf("expected 1 http request, got %d", snapshot.HTTPRequestsTotal)
	}
	if snapshot.GetKeyTotal != 1 {
		t.Fatalf("expected 1 get key, got %d", snapshot.GetKeyTotal)
	}
	if snapshot.SuccessTotal != 1 {
		t.Fatalf("expected 1 success, got %d", snapshot.SuccessTotal)
	}
}

func TestHandleKMIPGetKeyNotFound(t *testing.T) {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()

	dispatcher := usecase.NewDispatcher(repo, auditLogger)
	handler := HandleKMIP(dispatcher, metricsCollector)

	operationValue := make([]byte, 4)
	binary.BigEndian.PutUint32(operationValue, uint32(ttlv.OperationGet))

	body := []byte{
		0x42, 0x00, 0x5C,
		0x05,
		0x00, 0x00, 0x00, 0x04,
	}
	body = append(body, operationValue...)

	keyIDValue := []byte("missing-key")
	keyIDBlock := []byte{
		0x42, 0x00, 0x94,
		0x07,
		0x00, 0x00, 0x00, byte(len(keyIDValue)),
	}
	keyIDBlock = append(keyIDBlock, keyIDValue...)
	body = append(body, keyIDBlock...)

	req := httptest.NewRequest(http.MethodPost, "/kmip", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d, body: %s", rec.Code, rec.Body.String())
	}

	snapshot := metricsCollector.Snapshot(req.Context())
	if snapshot.HTTPRequestsTotal != 1 {
		t.Fatalf("expected 1 http request, got %d", snapshot.HTTPRequestsTotal)
	}
	if snapshot.GetKeyTotal != 1 {
		t.Fatalf("expected 1 get key, got %d", snapshot.GetKeyTotal)
	}
	if snapshot.NotFoundTotal != 1 {
		t.Fatalf("expected 1 not found, got %d", snapshot.NotFoundTotal)
	}
}

func TestHandleKMIPInvalidTTLV(t *testing.T) {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()

	dispatcher := usecase.NewDispatcher(repo, auditLogger)
	handler := HandleKMIP(dispatcher, metricsCollector)

	body := []byte{0x42, 0x00, 0x5C}

	req := httptest.NewRequest(http.MethodPost, "/kmip", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	snapshot := metricsCollector.Snapshot(req.Context())
	if snapshot.HTTPRequestsTotal != 1 {
		t.Fatalf("expected 1 http request, got %d", snapshot.HTTPRequestsTotal)
	}
	if snapshot.HTTPErrorsTotal != 1 {
		t.Fatalf("expected 1 http error, got %d", snapshot.HTTPErrorsTotal)
	}
}