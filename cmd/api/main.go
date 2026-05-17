package main

import (
	"log"
	"net/http"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/transport/httpapi"
	"kmipDemo/internal/usecase"
)

func main() {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()

	dispatcher := usecase.NewDispatcher(repo, auditLogger)

	mux := http.NewServeMux()

	mux.HandleFunc("/kmip", httpapi.HandleKMIP(dispatcher, metricsCollector))
	mux.HandleFunc("/metrics", httpapi.HandleMetrics(metricsCollector))
	mux.HandleFunc("/audit", httpapi.HandleAudit(auditLogger))
	mux.HandleFunc("/dashboard", httpapi.HandleDashboard())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"kmip-demo"}`))
	})

	addr := ":8080"

	log.Printf("starting kmip demo api on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}