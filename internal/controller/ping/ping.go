package ping

import (
	"database/sql"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

// TODO: Нужно получать значение dsn из the main service of the application- структура=служба со всеми основными сущностями нашего приложения
// См. services.Shortener - сделать имено тут

func HealthCheck(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// if config.FlagDsn == "" {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		//TODO: Получать из вызова, а не напрямую!
		// 1. Подключение к базе
		db, err := sql.Open("postgres", "postgres://postgres:qwerty@localhost:5434/postgres?sslmode=disable")
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
