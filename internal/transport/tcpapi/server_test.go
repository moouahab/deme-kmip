package tcpapi

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"testing"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase"
)

func TestHandleMessageCreateKey(t *testing.T) {
	server, repo := newTestServer()
	var out bytes.Buffer

	err := server.HandleMessage(context.Background(), bytes.NewReader(createKeyRequest()), &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var resp Response
	if err := json.NewDecoder(&out).Decode(&resp); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	if !resp.OK {
		t.Fatalf("expected ok response, got %+v", resp)
	}

	keys, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("expected no error while listing keys, got %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestHandleMessageInvalidTTLV(t *testing.T) {
	server, _ := newTestServer()
	var out bytes.Buffer

	err := server.HandleMessage(context.Background(), bytes.NewReader([]byte{0x42, 0x00}), &out)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var resp Response
	if err := json.NewDecoder(&out).Decode(&resp); err != nil {
		t.Fatalf("expected valid json response, got %v", err)
	}
	if resp.OK {
		t.Fatalf("expected error response, got %+v", resp)
	}
	if resp.Error != "bad_request" {
		t.Fatalf("expected bad_request, got %s", resp.Error)
	}
}

func newTestServer() (*Server, *kms.MemoryRepository) {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	collector := metrics.NewMemoryCollector()
	dispatcher := usecase.NewDispatcher(repo, auditLogger)

	return NewServer(dispatcher, collector), repo
}

func createKeyRequest() []byte {
	body := encodeEnumeration(ttlv.TagOperation, uint32(ttlv.OperationCreate))
	body = append(body, encodeEnumeration(ttlv.TagObjectType, uint32(ttlv.ObjectTypeSymmetricKey))...)
	return body
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
