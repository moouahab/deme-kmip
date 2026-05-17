package handlers

import (
	"context"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/usecase/models"
)

func GetKey(repo kms.Repository, auditLogger audit.Logger) func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
	return func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
		key, err := repo.Get(ctx, req.KeyID)
		if err != nil {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "get_key",
				KeyID:     req.KeyID,
				Result:    "not_found",
				Error:     err.Error(),
			})
			return models.OperationResponse{}, err
		}

		if key.Status == kms.KeyStatusDestroyed {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "get_key",
				KeyID:     req.KeyID,
				Status:    string(key.Status),
				Result:    "not_found",
				Error:     kms.ErrKeyNotFound.Error(),
			})
			return models.OperationResponse{}, kms.ErrKeyNotFound
		}

		_ = auditLogger.Log(ctx, audit.Event{
			Operation: "get_key",
			KeyID:     key.ID,
			Status:    string(key.Status),
			Result:    "success",
		})

		return models.OperationResponse{
			KeyID:  key.ID,
			Status: string(key.Status),
		}, nil
	}
}