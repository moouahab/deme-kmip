package httpapi

import (
	"net/http"

	"kmipDemo/internal/audit"
)

type AuditReader interface {
	Events() []audit.Event
}

func HandleAudit(reader AuditReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		events := reader.Events()

		if err := writeJSON(w, http.StatusOK, events); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "cannot write audit response")
			return
		}
	}
}
