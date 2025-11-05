package save

import (
	"app/internal/entity"
	"app/internal/repository/pg"
	"app/internal/usecase/random"
	"io"
	"log/slog"
	"net/http"
)

func New(log *slog.Logger, repo *pg.PostgresRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		shortUrl := entity.ShortURL{
			OriginalURL: string(body),
			ID:          random.NewRandomString(),
		}

		err = pg.CreateRecord(log, shortUrl, repo)
		if err != nil {
			log.Error("failed to add record")
			return
		}
	}
}
