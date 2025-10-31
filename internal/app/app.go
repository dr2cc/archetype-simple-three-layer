package app

import (
	"app/internal/config"
	"app/internal/repository/pg"
	"app/internal/usecase/logger/sl"
	myLog "app/internal/usecase/middleware/logger"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// - 2️⃣ Внедрение зависимостей.
// 		Инициализируйте и соедините (это и есть внедрение зависимостей)
// 		основные компоненты вашего приложения, такие как:
//  -- клиент базы данных,
//  -- уровень хранения (это компонент приложения, отвечающий за абстрагирование
// и управление взаимодействием с источником данных)
//  -- обработчики запросов (здесь они Application or Busines Logic).

// - 3️⃣ Настройка маршрутизатора: создайте экземпляр маршрутизатора Chi и зарегистрируйте маршруты,
// передав обработчики из логического уровня вашего приложения.

// - 4️⃣ Запуск сервера: запустите HTTP-сервер, обычно с помощью http.ListenAndServe, и корректно обработайте возможные ошибки запуска.

// Run creates objects via constructors.
func Run(cfg *config.Config) { //nolint: gocyclo,cyclop,funlen,gocritic,nolintlint
	log := setupLogger(cfg.Env)
	//log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении
	log.Info("init server", slog.String("address", cfg.HTTPServer.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	// Repository🧹🏦
	//pg.InitDB(log, cfg)
	db, err := pg.InitDB(log, cfg)
	if err != nil {
		log.Error("failed to connect storage")
		os.Exit(1)
	}
	// создаем/проверяем наличие таблицы
	errStorage := pg.New(log, db.DB)
	if errStorage != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}

	// // ...

	// Use-Case🧹🏦
	// ...

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	// Waiting signal🧹🏦
	// Логика Graceful Shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// HTTP Server🧹🏦
	router := chi.NewRouter()
	// Middleware встроенный в chi
	router.Use(middleware.RequestID) // Трассировка. Добавляется request_id в каждый запрос
	router.Use(middleware.Logger)    // Логирование всех запросов
	// Если внутри произойдет паника, приложение не упадет.
	// Recoverer это compress.Gzipper, которое восстанавливается после паники,
	// регистрирует панику и выводит идентификатор запроса, если он указан.
	router.Use(middleware.Recoverer)
	router.Use(myLog.New(log))       // Меняю логгер на мой
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов. Удалит суффикс из пути маршрутизации и продолжит маршрутизацию
	// Server startup parameters:
	httpServer := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
		//ReadTimeout:  cfg.HTTPServer.Timeout,
		//WriteTimeout: cfg.HTTPServer.Timeout,
		//IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Логика web-сервера
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	// Смысл таймаута был, но сейчас потерян..
	ctx := context.Background() //context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	// Shutdown🧹🏦
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))

		return
	}

	// TODO: close storage

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
