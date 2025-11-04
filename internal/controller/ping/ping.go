package ping

import (
	"app/internal/config"
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

func HealthCheck(log *slog.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// if config.FlagDsn == "" {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		// 1. Подключение к базе
		db, err := sql.Open("postgres", cfg.DSN)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//fmt.Fprint(w, "Error connecting to the database:", err)
			return
		}
		//w.WriteHeader(http.StatusOK)
		//fmt.Fprint(w, dn, " - successfully connected to the database!")
		defer db.Close()
	}
}
