package db

import (
	"context"
	"log"
	"time"

	"github.com/btynybekov/marketplace/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool — глобальный пул подключений
var Pool *pgxpool.Pool

// InitDB инициализирует подключение к Postgres через pgxpool
func InitDB(conf config.EnvConfig) {
	dsn := config.GetDSN(conf)
	log.Println("Connecting to DB...")

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DATABASE_URL: %v", err)
	}

	// Настройка пула
	cfg.MaxConns = 20
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.HealthCheckPeriod = 5 * time.Minute

	// Создаём пул
	Pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to create DB pool: %v", err)
	}

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := Pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("✅ Database connected successfully via pgxpool!")
}

// CloseDB закрывает пул соединений
func CloseDB() {
	if Pool != nil {
		Pool.Close()
		log.Println("Database connection closed")
	}
}
