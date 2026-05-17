package ttlv

import (
	"encoding/binary"
	"fmt"
)

func DecodeBlocks(data []byte) ([]Block, error) {
	var blocks []Block
	offset := 0

	for offset < len(data) {
		block, consumed, err := decodeBlock(data[offset:])
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
		offset += consumed
	}
	return blocks, nil
}

func decodeBlock(data []byte) (Block, int, error) {
	
	if len(data) < HeaderSize {
		return Block{}, 0, fmt.Errorf("ttlv: message too short: got %d bytes, need at least %d", len(data), HeaderSize)
	}
	
	tag := Tag(uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2]))
	ttlvType := Type(data[3])
	length := binary.BigEndian.Uint32(data[4:8])
	
	if length > MaxValueLength {
		return Block{}, 0, fmt.Errorf("ttlv: value too large: %d bytes", length)
	}
	
	totalSize := HeaderSize + int(length)
	if len(data) < totalSize {
		return Block{}, 0, fmt.Errorf("ttlv: invalid length: expected %d value bytes, got %d", length, len(data)-HeaderSize)
	}
	
	value := data[HeaderSize:totalSize]
	block := Block{
		Tag:    tag,
		Type:   ttlvType,
		Length: length,
		Value:  value,
	}
	if err := block.Validate(); err != nil {
		return Block{}, 0, err
	}
	
	return block, totalSize, nil
}
