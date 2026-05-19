package usecase

import (
	"context"
	"fmt"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/usecase/handlers"
	"kmipDemo/internal/usecase/models"
)

func NewDispatcher(repo kms.Repository, auditLogger audit.Logger) *Dispatcher {
	d := &Dispatcher{
		repo:        repo,
		auditLogger: auditLogger,
	}

	d.handlers = map[ttlv.Operation]OperationHandler{
		ttlv.OperationCreate:        handlers.CreateKey(repo, auditLogger),
		ttlv.OperationGet:           handlers.GetKey(repo, auditLogger),
		ttlv.OperationDestroy:       handlers.DestroyKey(repo, auditLogger),
		ttlv.OperationActivate:      handlers.ActivateKey(repo, auditLogger),
		ttlv.OperationRevoke:        handlers.RevokeKey(repo, auditLogger),
		ttlv.OperationLocate:        handlers.LocateKeys(repo, auditLogger),
		ttlv.OperationGetAttributes: handlers.GetKeyAttributes(repo, auditLogger),
	}

	return d
}

func (d *Dispatcher) Dispatch(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
	handler, ok := d.handlers[req.Operation]
	if !ok {
		return models.OperationResponse{}, fmt.Errorf("usecase: unsupported operation: %d", req.Operation)
	}

	return handler(ctx, req)
}
