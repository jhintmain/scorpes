package api

import (
	"fmt"
	"net/http"

	"github.com/hooneun/scorpes/internal/config"
	db "github.com/hooneun/scorpes/internal/db/sqlc"
)

type API struct {
	cfg *config.Config
	db  *db.Queries
}

func NewAPI(cfg *config.Config, db *db.Queries) *API {
	return &API{
		cfg: cfg,
		db:  db,
	}
}

func (a *API) setupRouter() http.Handler {
	r := NewRouter()

	RegisterRoutes(r)

	return r
}

func (a *API) Run() error {
	addr := fmt.Sprintf(":%d", 8090)
	server := &http.Server{
		Addr:    addr,
		Handler: a.setupRouter(),
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
