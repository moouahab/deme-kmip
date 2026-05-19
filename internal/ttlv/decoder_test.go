package ttlv

import (
	"encoding/binary"
	"testing"
)

func TestDecodeBlocksValidOperation(t *testing.T) {
	data := []byte{
		0x42, 0x00, 0x5C, // TagOperation
		0x05,                   // TypeEnumeration
		0x00, 0x00, 0x00, 0x04, // Length = 4
		0x00, 0x00, 0x00, 0x01, // OperationCreate
	}

	blocks, err := DecodeBlocks(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]

	if block.Tag != TagOperation {
		t.Fatalf("expected tag %X, got %X", TagOperation, block.Tag)
	}

	if block.Type != TypeEnumeration {
		t.Fatalf("expected type %X, got %X", TypeEnumeration, block.Type)
	}

	if block.Length != 4 {
		t.Fatalf("expected length 4, got %d", block.Length)
	}

	got := Operation(binary.BigEndian.Uint32(block.Value))
	if got != OperationCreate {
		t.Fatalf("expected operation create, got %v", got)
	}
}

func TestDecodeBlocksMessageTooShort(t *testing.T) {
	data := []byte{0x42, 0x00, 0x5C}

	_, err := DecodeBlocks(data)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDecodeBlocksInvalidLength(t *testing.T) {
	data := []byte{
		0x42, 0x00, 0x5C, // TagOperation
		0x05,                   // TypeEnumeration
		0x00, 0x00, 0x00, 0x04, // Length = 4
		0x00, 0x00, // only 2 bytes instead of 4
	}

	_, err := DecodeBlocks(data)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
