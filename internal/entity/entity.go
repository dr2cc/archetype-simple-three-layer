// Package models contains business models description.
package entity

import (
	"app/internal/config"
	"app/internal/repository/pg"
	"log/slog"

	"github.com/go-chi/chi"
)

// Пока предположение, что правильно не знаю
type Entity struct {
	Router *chi.Mux
	Log    *slog.Logger
	Cfg    *config.Config
	DB     *pg.Postgres
}
