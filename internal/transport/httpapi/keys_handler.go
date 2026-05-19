package httpapi

import (
	"context"
	"net/http"

	"kmipDemo/internal/kms"
)

type KeyLister interface {
	List(ctx context.Context) ([]kms.Key, error)
}

func HandleKeys(lister KeyLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		keys, err := lister.List(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "cannot list keys")
			return
		}

		if err := writeJSON(w, http.StatusOK, keys); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "cannot write keys response")
			return
		}
	}
}
