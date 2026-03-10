package main

import (
	"context"
	"log"

	"github.com/hooneun/scorpes/internal/api"
	"github.com/hooneun/scorpes/internal/config"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
	"github.com/hooneun/scorpes/internal/scheduler"
	"github.com/hooneun/scorpes/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	// -------------------------
	// PostgreSQL connection
	// -------------------------
	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.Server.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	queries := db.New(dbPool)
	//

	a := api.NewAPI(cfg, queries)

	pool := worker.NewPool(5, 100, cfg, queries)
	pool.Start()

	cronScheduler := scheduler.NewCronScheduler(pool.JobQueue, cfg, queries)
	cronScheduler.Start()

	if err := a.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
