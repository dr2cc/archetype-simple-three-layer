package pg

import (
	"app/internal/config"
	"app/internal/usecase/logger/sl"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	// Без такой конструкции не возможно подключится к postgres
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sql.DB
}

// Инициализация подключения к PostgreSQL
func InitDB(log *slog.Logger, cfg *config.Config) (*Postgres, error) {
	// Getting DSN from environment variables
	//dsn := os.Getenv("DATABASE_DSN")

	// // dsn проверяем перед вызовом InitDB
	// if dsn == "" {
	// 	log.Error("DATABASE_DSN not specified in env")
	// 	//os.Exit(1)
	// }

	// 1. Подключение к базе
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Error("DB connection error", sl.Err(err))
		return nil, fmt.Errorf("connection error: %v", err)
	}

	// // Не забыть про defer!!
	//defer db.Close()

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяю подключение с таймаутом ответа
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// освобождаем ресурс
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Error("error to ping", sl.Err(err))
		return nil, fmt.Errorf("error to ping db: %v", err)
	}

	return &Postgres{DB: db}, nil
}

func New(log *slog.Logger, db *sql.DB) error {

	// 2. Создаем таблицу, если ее еще нет
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS aliases(
        alias VARCHAR NOT NULL UNIQUE,
        url TEXT NOT NULL);
	`)

	if err != nil {
		log.Error(err.Error())
	}

	// Отправляем комманду (CREATE TABLE в данном случае)
	// Exec выполняет подготовленный оператор (stmt) с заданными аргументами
	// и возвращает [Result], суммирующий эффект оператора.
	// В данной ситуации этот "эффект" не используется
	_, err = stmt.Exec()
	if err != nil {
		log.Error(err.Error())
	}

	return nil
}
