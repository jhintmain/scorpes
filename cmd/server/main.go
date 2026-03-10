package main

import (
	"context"
	"log"

	"github.com/hooneun/scorpes/internal/api"
	"github.com/hooneun/scorpes/internal/config"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
	"github.com/hooneun/scorpes/internal/scheduler"
	"github.com/hooneun/scorpes/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	queries := db.New(pool)

	a := api.NewAPI(cfg, queries)

	log.Printf("Server starting on port %s", cfg.Server.Port)

	if err := a.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
