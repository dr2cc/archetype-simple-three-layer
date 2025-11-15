package server

import (
	"app/internal/config"
	v1 "app/internal/controller/http/v1"
	"app/internal/repository/pg"
	"app/internal/usecase/logger/sl"
	"app/internal/usecase/random"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
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

// В целом не до конца понимаю, что это дает (04.11.2025)
// но давно хотел создать в конструкторе "главную" структуру приложения
type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	return &App{}
}

// Run creates objects (via constructors!)
func (a *App) Run(cfg *config.Config) {
	log := setupLogger(cfg.Env)
	log.Info("init server", slog.String("address", cfg.HTTPServer.Address)) // Помимо сообщения выведем параметр с адресом
	log.Debug("logger debug mode enabled")

	// Repository🧹🏦
	repo, err := pg.NewPostgresRepo(log, cfg)
	if err != nil {
		log.Error("failed to connect storage")
		os.Exit(1)
	}

	// Use-Case🧹🏦
	// В данный момент (14.11.25) именно service здесь я не создаю
	randomKey := random.RandomGenerator{}

	// Router
	mux := chi.NewRouter()
	// middlewares & handlers
	v1.Router(mux, cfg, repo, randomKey, log)

	// ❗Graceful shutdown
	// done: Это наш "стоп-кран".
	// Это буферизованный канал, который будет ожидать системные сигналы.
	done := make(chan os.Signal, 1)
	// signal.Notify: Регистрирует канал done для получения уведомлений,
	// когда операционная система отправляет сигналы прерывания
	// (Ctrl+C), SIGINT или SIGTERM
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// HTTP Server - конфигурация и запуск
	a.httpServer = &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Отдельная горутина: сервер запускается в своей собственной горутине.
	// Это необходимо, так как ListenAndServe() является блокирующим вызовом.
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	// Ожидание сигнала остановки.
	// <-done: Это критическая точка синхронизации. Основная горутина main блокируется здесь.
	// Она будет ждать, пока в канал done не придет системный сигнал.
	// Как только пользователь нажимает Ctrl+C, канал разблокируется, и выполнение продолжается.
	<-done
	log.Info("stopping server")

	// Корректное завершение с таймаутом (context.WithTimeout и Shutdown).
	// context.WithTimeout: создает контекст, который автоматически отменится через 10 секунд.
	// Это "страховка" от зависания сервера.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	// Всегда отменяем контекст, чтобы освободить его ресурсы
	defer cancel()

	// srv.Shutdown(ctx): вызывает изящное (graceful) завершение работы.
	// Он перестает принимать новые запросы, но дает активным запросам время завершиться.
	// Он использует канал <-ctx.Done() (который находится внутри ctx), чтобы узнать, когда истечет 10-секундный лимит.
	if err := a.httpServer.Shutdown(ctx); err != nil {
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
