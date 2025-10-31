package ping

//_ "github.com/lib/pq"

// TODO: Нужно получать значение dsn из the main service of the application- структура=служба со всеми основными сущностями нашего приложения
// См. services.Shortener - сделать имено тут

// func HealthCheck(log *slog.Logger) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		if config.FlagDsn == "" {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		// 1. Подключение к базе
// 		db, err := sql.Open("postgres", config.FlagDsn)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		err = db.Ping()
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			//fmt.Fprint(w, "Error connecting to the database:", err)
// 			return
// 		}
// 		//w.WriteHeader(http.StatusOK)
// 		//fmt.Fprint(w, dn, " - successfully connected to the database!")
// 		defer db.Close()
// 	}
// }
