package save

import (
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

		url := string(body)
		alias := random.NewRandomString()
		// TODO: здесь создать объект Entry (в конструкторе?!) типа entity.ShortURL,
		// в него передать url и alias , а затем работать с ними
		err = pg.CreateRecord(log, url, alias, repo)
		if err != nil {
			log.Error("failed to add record")
			return
		}
	}
}
