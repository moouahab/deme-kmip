package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleKMIPInvalidTTLV(t *testing.T) {
	f := newHTTPFixture()
	rec := serveKMIP(t, f, []byte{0x42, 0x00, 0x5C})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}

	snapshot := f.metrics.Snapshot(t.Context())
	if snapshot.HTTPRequestsTotal != 1 || snapshot.HTTPErrorsTotal != 1 {
		t.Fatalf("unexpected metrics snapshot: %+v", snapshot)
	}
}

func TestHandleKMIPBodyTooLarge(t *testing.T) {
	f := newHTTPFixture()
	body := bytes.Repeat([]byte{0x00}, maxKMIPRequestBodySize+1)
	rec := serveKMIP(t, f, body)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", rec.Code)
	}
}

func TestHandleKMIPRejectsWrongMethod(t *testing.T) {
	f := newHTTPFixture()
	req := httptest.NewRequest(http.MethodGet, "/kmip", nil)
	rec := httptest.NewRecorder()

	f.handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}
