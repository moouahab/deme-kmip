package handlers

import (
	"context"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/usecase/models"
)

func RevokeKey(repo kms.Repository, auditLogger audit.Logger) func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
	return func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
		key, err := repo.Get(ctx, req.KeyID)
		if err != nil {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "revoke_key",
				KeyID:     req.KeyID,
				Result:    "not_found",
				Error:     err.Error(),
			})
			return models.OperationResponse{}, err
		}

		if key.Status == kms.KeyStatusDestroyed {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "revoke_key",
				KeyID:     req.KeyID,
				Status:    string(key.Status),
				Result:    "not_found",
				Error:     kms.ErrKeyNotFound.Error(),
			})
			return models.OperationResponse{}, kms.ErrKeyNotFound
		}

		key.Status = kms.KeyStatusRevoked
		updated, err := repo.Update(ctx, key)
		if err != nil {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "revoke_key",
				KeyID:     req.KeyID,
				Result:    "error",
				Error:     err.Error(),
			})
			return models.OperationResponse{}, err
		}

		_ = auditLogger.Log(ctx, audit.Event{
			Operation: "revoke_key",
			KeyID:     updated.ID,
			Status:    string(updated.Status),
			Result:    "success",
		})

		return models.OperationResponse{
			KeyID:  updated.ID,
			Status: string(updated.Status),
		}, nil
	}
}
