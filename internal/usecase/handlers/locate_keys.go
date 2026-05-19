package handlers

import (
	"context"
	"time"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/usecase/models"
)

func LocateKeys(repo kms.Repository, auditLogger audit.Logger) func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
	return func(ctx context.Context, req models.OperationRequest) (models.OperationResponse, error) {
		keys, err := repo.List(ctx)
		if err != nil {
			_ = auditLogger.Log(ctx, audit.Event{
				Operation: "locate_keys",
				Result:    "error",
				Error:     err.Error(),
			})
			return models.OperationResponse{}, err
		}

		summaries := make([]models.KeySummary, 0, len(keys))
		for _, key := range keys {
			if key.Status == kms.KeyStatusDestroyed {
				continue
			}
			summaries = append(summaries, keySummary(key))
		}

		_ = auditLogger.Log(ctx, audit.Event{
			Operation: "locate_keys",
			Result:    "success",
		})

		return models.OperationResponse{
			Keys: summaries,
		}, nil
	}
}

func keySummary(key kms.Key) models.KeySummary {
	return models.KeySummary{
		ID:         key.ID,
		ObjectType: uint32(key.ObjectType),
		Status:     string(key.Status),
		CreatedAt:  formatTime(key.CreatedAt),
		UpdatedAt:  formatTime(key.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339Nano)
}
