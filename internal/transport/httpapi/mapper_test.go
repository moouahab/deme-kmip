package httpapi

import (
	"encoding/binary"
	"testing"

	"kmipDemo/internal/ttlv"
)

func TestBlocksToOperationRequestCreateKey(t *testing.T) {
	operationValue := make([]byte, 4)
	binary.BigEndian.PutUint32(operationValue, uint32(ttlv.OperationCreate))

	objectTypeValue := make([]byte, 4)
	binary.BigEndian.PutUint32(objectTypeValue, uint32(ttlv.ObjectTypeSymmetricKey))

	blocks := []ttlv.Block{
		{
			Tag:    ttlv.TagOperation,
			Type:   ttlv.TypeEnumeration,
			Length: 4,
			Value:  operationValue,
		},
		{
			Tag:    ttlv.TagObjectType,
			Type:   ttlv.TypeEnumeration,
			Length: 4,
			Value:  objectTypeValue,
		},
	}

	req, err := blocksToOperationRequest(blocks)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Operation != ttlv.OperationCreate {
		t.Fatalf("expected OperationCreate, got %v", req.Operation)
	}

	if req.ObjectType != ttlv.ObjectTypeSymmetricKey {
		t.Fatalf("expected ObjectTypeSymmetricKey, got %v", req.ObjectType)
	}
}

func TestBlocksToOperationRequestGetKey(t *testing.T) {
	operationValue := make([]byte, 4)
	binary.BigEndian.PutUint32(operationValue, uint32(ttlv.OperationGet))

	blocks := []ttlv.Block{
		{
			Tag:    ttlv.TagOperation,
			Type:   ttlv.TypeEnumeration,
			Length: 4,
			Value:  operationValue,
		},
		{
			Tag:    ttlv.TagUniqueIdentifier,
			Type:   ttlv.TypeTextString,
			Length: uint32(len("key-123")),
			Value:  []byte("key-123"),
		},
	}

	req, err := blocksToOperationRequest(blocks)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Operation != ttlv.OperationGet {
		t.Fatalf("expected OperationGet, got %v", req.Operation)
	}

	if req.KeyID != "key-123" {
		t.Fatalf("expected key-123, got %s", req.KeyID)
	}
}

func TestBlocksToOperationRequestInvalidOperationType(t *testing.T) {
	blocks := []ttlv.Block{
		{
			Tag:    ttlv.TagOperation,
			Type:   ttlv.TypeTextString,
			Length: 4,
			Value:  []byte{0x00, 0x00, 0x00, 0x01},
		},
	}

	_, err := blocksToOperationRequest(blocks)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBlocksToOperationRequestInvalidOperationLength(t *testing.T) {
	blocks := []ttlv.Block{
		{
			Tag:    ttlv.TagOperation,
			Type:   ttlv.TypeEnumeration,
			Length: 2,
			Value:  []byte{0x00, 0x01},
		},
	}

	_, err := blocksToOperationRequest(blocks)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBlocksToOperationRequestMissingOperation(t *testing.T) {
	blocks := []ttlv.Block{
		{
			Tag:    ttlv.TagUniqueIdentifier,
			Type:   ttlv.TypeTextString,
			Length: uint32(len("key-123")),
			Value:  []byte("key-123"),
		},
	}

	_, err := blocksToOperationRequest(blocks)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
