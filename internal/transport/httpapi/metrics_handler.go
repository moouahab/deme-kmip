package httpapi

import (
	"net/http"

	"kmipDemo/internal/metrics"
)

func HandleMetrics(collector metrics.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}

		snapshot := collector.Snapshot(r.Context())

		if err := writeJSON(w, http.StatusOK, snapshot); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", "cannot write metrics response")
			return
		}
	}
}