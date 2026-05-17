package usecase

import (
	"context"
	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/models"
)

type OperationResponse struct {
	KeyID  string
	Status string
}

type OperationHandler func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error)

type Dispatcher struct {
	repo        kms.Repository
	auditLogger audit.Logger
	handlers    map[ttlv.Operation]OperationHandler
}