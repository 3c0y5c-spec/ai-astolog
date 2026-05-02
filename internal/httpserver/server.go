package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	addr   string
	logger *slog.Logger
	server *http.Server
}

func New(addr string, logger *slog.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handleHealthz)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		addr:   addr,
		logger: logger,
		server: server,
	}
}

func (s *Server) ListenAndServe() error {
	s.logger.Info("http server started", "addr", s.addr)
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("http server shutting down")
	return s.server.Shutdown(ctx)
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
