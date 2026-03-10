package api

import (
	"fmt"
	"net/http"

	"github.com/hooneun/scorpes/internal/config"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
)

type API struct {
	cfg     *config.Config
	queries *db.Queries
}

func NewAPI(cfg *config.Config, queries *db.Queries) *API {
	return &API{
		cfg:     cfg,
		queries: queries,
	}
}

func (a *API) setupRouter() http.Handler {
	r := NewRouter()

	targetHandler := NewTargetHandler(a.queries)
	RegisterRoutes(r, targetHandler)

	return r
}

func (a *API) Run() error {
	addr := fmt.Sprintf(":%s", a.cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: a.setupRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
