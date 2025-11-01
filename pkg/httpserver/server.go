// Package httpserver implements HTTP server.
package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	_defaultReadTimeout  = 5 * time.Second
	_defaultWriteTimeout = 5 * time.Second
)

// New -.
func New(addr string, router *chi.Mux, log *slog.Logger) *http.Server {
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	// Логика web-сервера
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	return httpServer
}
