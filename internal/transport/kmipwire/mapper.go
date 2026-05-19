package kmipwire

import (
	"encoding/binary"
	"fmt"

	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

func BlocksToOperationRequest(blocks []ttlv.Block) (models.OperationRequest, error) {
	var req models.OperationRequest

	for _, block := range blocks {
		switch block.Tag {
		case ttlv.TagOperation:
			if block.Type != ttlv.TypeEnumeration {
				return models.OperationRequest{}, fmt.Errorf("operation must be enumeration")
			}
			if len(block.Value) != 4 {
				return models.OperationRequest{}, fmt.Errorf("operation value must be 4 bytes")
			}
			req.Operation = ttlv.Operation(binary.BigEndian.Uint32(block.Value))

		case ttlv.TagUniqueIdentifier:
			if block.Type != ttlv.TypeTextString {
				return models.OperationRequest{}, fmt.Errorf("key id must be text string")
			}
			req.KeyID = string(block.Value)

		case ttlv.TagObjectType:
			if block.Type != ttlv.TypeEnumeration {
				return models.OperationRequest{}, fmt.Errorf("object type must be enumeration")
			}
			if len(block.Value) != 4 {
				return models.OperationRequest{}, fmt.Errorf("object type value must be 4 bytes")
			}
			req.ObjectType = ttlv.ObjectType(binary.BigEndian.Uint32(block.Value))

		case ttlv.TagRequestPayload:
			req.Payload = block.Value
		}
	}

	if !req.Operation.IsValid() {
		return models.OperationRequest{}, fmt.Errorf("invalid or missing operation")
	}
	switch req.Operation {
	case ttlv.OperationCreate:
		if !req.ObjectType.IsValid() {
			return models.OperationRequest{}, fmt.Errorf("invalid or missing object type")
		}
	case ttlv.OperationGet, ttlv.OperationDestroy, ttlv.OperationActivate, ttlv.OperationRevoke, ttlv.OperationGetAttributes:
		if req.KeyID == "" {
			return models.OperationRequest{}, fmt.Errorf("missing key id")
		}
	}

	return req, nil
}
