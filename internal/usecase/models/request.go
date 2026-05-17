package models

import "kmipDemo/internal/ttlv"

type OperationRequest struct {
	Operation  ttlv.Operation
	KeyID      string
	ObjectType ttlv.ObjectType
	Payload    []byte
}
