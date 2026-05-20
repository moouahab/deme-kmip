package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"kmipDemo/internal/audit"
	"kmipDemo/internal/kms"
	"kmipDemo/internal/metrics"
	"kmipDemo/internal/transport/httpapi"
	"kmipDemo/internal/transport/tcpapi"
	"kmipDemo/internal/usecase"
)

func main() {
	repo := kms.NewMemoryRepository()
	auditLogger := audit.NewMemoryLogger()
	metricsCollector := metrics.NewMemoryCollector()

	dispatcher := usecase.NewDispatcher(repo, auditLogger)

	mux := http.NewServeMux()

	mux.HandleFunc("/kmip", httpapi.HandleKMIP(dispatcher, metricsCollector))
	mux.HandleFunc("/keys", httpapi.HandleKeys(repo))
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

	go func() {
		tcpServer := tcpapi.NewServer(dispatcher, metricsCollector)
		if err := tcpServer.ListenAndServe(context.Background(), ":5696"); err != nil {
			log.Printf("tcp transport stopped: %v", err)
		}
	}()

	addr := ":8080"
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("starting kmip demo api on %s", addr)
	
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
