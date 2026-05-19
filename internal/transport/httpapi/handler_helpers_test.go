package httpapi

import (
	"bytes"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase"
)

type httpFixture struct {
	repo      *kms.MemoryRepository
	metrics   *metrics.MemoryCollector
	handler   http.HandlerFunc
	auditLogs *audit.MemoryLogger
}

func newHTTPFixture() httpFixture {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()
	dispatcher := usecase.NewDispatcher(repo, auditLogger)

	return httpFixture{
		repo:      repo,
		metrics:   metricsCollector,
		handler:   HandleKMIP(dispatcher, metricsCollector),
		auditLogs: auditLogger,
	}
}

func serveKMIP(t *testing.T, f httpFixture, body []byte) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/kmip", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	f.handler.ServeHTTP(rec, req)
	return rec
}

func encodeEnumeration(tag ttlv.Tag, value uint32) []byte {
	valueBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(valueBytes, value)

	body := []byte{
		byte(uint32(tag) >> 16),
		byte(uint32(tag) >> 8),
		byte(uint32(tag)),
		byte(ttlv.TypeEnumeration),
		0x00, 0x00, 0x00, 0x04,
	}
	return append(body, valueBytes...)
}

func encodeText(tag ttlv.Tag, value string) []byte {
	valueBytes := []byte(value)
	body := []byte{
		byte(uint32(tag) >> 16),
		byte(uint32(tag) >> 8),
		byte(uint32(tag)),
		byte(ttlv.TypeTextString),
		0x00, 0x00, 0x00, byte(len(valueBytes)),
	}
	return append(body, valueBytes...)
}

func createKeyBody() []byte {
	body := encodeEnumeration(ttlv.TagOperation, uint32(ttlv.OperationCreate))
	body = append(body, encodeEnumeration(ttlv.TagObjectType, uint32(ttlv.ObjectTypeSymmetricKey))...)
	return body
}

func getKeyBody(keyID string) []byte {
	body := encodeEnumeration(ttlv.TagOperation, uint32(ttlv.OperationGet))
	body = append(body, encodeText(ttlv.TagUniqueIdentifier, keyID)...)
	return body
}
