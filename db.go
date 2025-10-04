package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Name     string
	SSLMode  string
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Host:    getenv("PGHOST", "localhost"),
		Port:    getenv("PGPORT", "5432"),
		User:    getenv("PGUSER", "user"),
		Pass:    getenv("PGPASSWORD", "password"),
		Name:    getenv("PGDATABASE", "rag_db"),
		SSLMode: getenv("PGSSLMODE", "disable"),
	}
}

func openPool(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := loadDBConfig()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	pcfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	// Reasonable defaults
	pcfg.MinConns = 0
	pcfg.MaxConns = 4
	pcfg.MaxConnLifetime = 30 * time.Minute
	return pgxpool.NewWithConfig(ctx, pcfg)
}
