package save

import (
	"app/internal/repository/pg"
	"app/internal/usecase/random"
	url_shortening_service "app/internal/usecase/urlshorteningservice"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	Repo *pg.PostgresRepo
}

func (h Handler) New(randomKey random.RandomGenerator, log *slog.Logger) http.HandlerFunc {
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

		shortUrl := url_shortening_service.NewService(string(url), randomKey)

		err = pg.NewRecord(log, shortUrl, h.Repo)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "failed to add record", http.StatusBadRequest)
			return
		}
	}
}
