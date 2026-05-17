package handlers

import (
	"context"

	"github.com/google/uuid"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/usecase/models"
)

func CreateKey(repo kms.Repository, auditLogger audit.Logger) func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
	return func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
		key := kms.Key{
			ID:         "key-" + uuid.NewString(),
			ObjectType: req.ObjectType,
			Status:     kms.KeyStatusActive,
		}

		created, err := repo.Create(ctx, key)
		if err != nil {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "create_key",
				KeyID:     key.ID,
				Result:    "error",
				Error:     err.Error(),
			})
			return models.OperationResponse{}, err
		}

		_ = auditLogger.Log(ctx, audit.Event{
			Operation: "create_key",
			KeyID:     created.ID,
			Status:    string(created.Status),
			Result:    "success",
		})

		return models.OperationResponse{
			KeyID:  created.ID,
			Status: string(created.Status),
		}, nil
	}
}