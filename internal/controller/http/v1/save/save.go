package save

import (
	"app/internal/repository/pg"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = pg.CreateRecord(log, string(body), db)
		if err != nil {
			log.Error("failed to add record")
			return
		}
	}
}
