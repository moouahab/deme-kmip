package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"kmipDemo/internal/kms"
	"kmipDemo/internal/ttlv"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/usecase"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return fmt.Errorf("httpapi: encode json response: %w", err)
	}

	return nil
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	if err := writeJSON(w, status, ErrorResponse{
		Error:   code,
		Message: message,
	}); err != nil {
		log.Printf("httpapi: failed to write error response: %v", err)
	}
}

func HandleKMIP(dispatcher *usecase.Dispatcher, collector metrics.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		collector.IncHTTPRequests(r.Context())
		if r.Method != http.MethodPost {
			collector.IncHTTPErrors(r.Context())
			writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
			return
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			collector.IncHTTPErrors(r.Context())
			writeError(w, http.StatusBadRequest, "bad_request", "cannot read request body")
			return
		}
		blocks, err := ttlv.DecodeBlocks(data)
		if err != nil {
			collector.IncHTTPErrors(r.Context())
			writeError(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}

		req, err := blocksToOperationRequest(blocks)
		if err != nil {
			collector.IncHTTPErrors(r.Context())
			writeError(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}
		switch req.Operation {
			case ttlv.OperationCreate:
				collector.IncCreateKey(r.Context())
			case ttlv.OperationGet:
				collector.IncGetKey(r.Context())
			case ttlv.OperationDestroy:
				collector.IncDestroyKey(r.Context())
		}
		resp, err := dispatcher.Dispatch(r.Context(), req)
		if err != nil {
			if errors.Is(err, kms.ErrKeyNotFound) {
				collector.IncNotFound(r.Context())
				writeError(w, http.StatusNotFound, "not_found", err.Error())
				return
			}
			collector.IncHTTPErrors(r.Context())
			writeError(w, http.StatusBadRequest, "bad_request", err.Error())
			return
		}
		collector.IncSuccess(r.Context())
		if err := writeJSON(w, http.StatusOK, resp); err != nil {
			log.Printf("httpapi: failed to write success response: %v", err)
		}
	}
}