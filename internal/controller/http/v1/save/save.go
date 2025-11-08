package save

import (
	"io"
	"log/slog"
	"net/http"

	"app/internal/entity"
	"app/internal/repository/pg"
	"app/internal/usecase/random"
)

type Handler struct {
	repo *pg.PostgresRepo
}

func New(repo *pg.PostgresRepo, randomKey random.RandomGenerator, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if string(url) == "" {
			http.Error(w, "content required", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		shortUrl := entity.ShortURL{
			OriginalURL: string(url),
			ID:          randomKey.NewRandomString(),
		}

		err = pg.CreateRecord(log, shortUrl, repo)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to add record", http.StatusBadRequest)
			return
		}
	}
}
