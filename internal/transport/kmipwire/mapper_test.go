package kmipwire

import (
	"encoding/binary"
	"testing"

	"kmipDemo/internal/ttlv"
)

func enumBlock(tag ttlv.Tag, value uint32) ttlv.Block {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, value)

	return ttlv.Block{
		Tag:    tag,
		Type:   ttlv.TypeEnumeration,
		Length: 4,
		Value:  bytes,
	}
}

func TestBlocksToOperationRequestCreateKey(t *testing.T) {
	req, err := BlocksToOperationRequest([]ttlv.Block{
		enumBlock(ttlv.TagOperation, uint32(ttlv.OperationCreate)),
		enumBlock(ttlv.TagObjectType, uint32(ttlv.ObjectTypeSymmetricKey)),
	})
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
	req, err := BlocksToOperationRequest([]ttlv.Block{
		enumBlock(ttlv.TagOperation, uint32(ttlv.OperationGet)),
		{
			Tag:    ttlv.TagUniqueIdentifier,
			Type:   ttlv.TypeTextString,
			Length: uint32(len("key-123")),
			Value:  []byte("key-123"),
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if req.KeyID != "key-123" {
		t.Fatalf("expected key-123, got %s", req.KeyID)
	}
}

func TestBlocksToOperationRequestRequiredFields(t *testing.T) {
	tests := []struct {
		name  string
		block ttlv.Block
	}{
		{name: "create requires object type", block: enumBlock(ttlv.TagOperation, uint32(ttlv.OperationCreate))},
		{name: "get requires key id", block: enumBlock(ttlv.TagOperation, uint32(ttlv.OperationGet))},
		{name: "get attributes requires key id", block: enumBlock(ttlv.TagOperation, uint32(ttlv.OperationGetAttributes))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BlocksToOperationRequest([]ttlv.Block{tt.block})
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestBlocksToOperationRequestLocateDoesNotRequireKeyID(t *testing.T) {
	req, err := BlocksToOperationRequest([]ttlv.Block{
		enumBlock(ttlv.TagOperation, uint32(ttlv.OperationLocate)),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if req.Operation != ttlv.OperationLocate {
		t.Fatalf("expected OperationLocate, got %v", req.Operation)
	}
}

func TestBlocksToOperationRequestInvalidOperation(t *testing.T) {
	_, err := BlocksToOperationRequest([]ttlv.Block{
		{
			Tag:    ttlv.TagOperation,
			Type:   ttlv.TypeTextString,
			Length: 4,
			Value:  []byte{0x00, 0x00, 0x00, 0x01},
		},
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
