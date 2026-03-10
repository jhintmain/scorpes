package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	db "github.com/hooneun/scorpes/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type TargetQuerier interface {
	ListTargets(ctx context.Context) ([]db.Target, error)
	CreateTarget(ctx context.Context, arg db.CreateTargetParams) (db.Target, error)
	UpdateTarget(ctx context.Context, arg db.UpdateTargetParams) (db.Target, error)
	GetTargetByID(ctx context.Context, id pgtype.UUID) (db.Target, error)
	SoftDeleteTarget(ctx context.Context, id pgtype.UUID) error
}

type TargetHandler struct {
	queries TargetQuerier
}

func NewTargetHandler(queries *db.Queries) *TargetHandler {
	return &TargetHandler{
		queries: queries,
	}
}

func NewTargetHandlerWithQuerier(queries TargetQuerier) *TargetHandler {
	return &TargetHandler{
		queries: queries,
	}
}

type CreateTargetRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Method          string `json:"method"`
	IntervalSeconds int32  `json:"interval_seconds"`
	TimeoutSeconds  int32  `json:"timeout_seconds"`
}

func (r *CreateTargetRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(r.URL) == "" {
		return errors.New("url is required")
	}

	if _, err := url.ParseRequestURI(r.URL); err != nil {
		return errors.New("invalid URL format")
	}

	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true,
		"DELETE": true, "HEAD": true, "PATCH": true, "OPTIONS": true,
	}

	method := strings.ToUpper(strings.TrimSpace(r.Method))
	if method == "" {
		r.Method = "GET"
	} else if !validMethods[method] {
		return errors.New("invalid HTTP method")
	} else {
		r.Method = method
	}

	if r.IntervalSeconds < 60 {
		return errors.New("interval seconds must be greater than or equal to 60")
	}

	if r.TimeoutSeconds < 1 {
		r.TimeoutSeconds = 10
	}

	return nil
}

func (h *TargetHandler) ListTargets(w http.ResponseWriter, r *http.Request) {
	targets, err := h.queries.ListTargets(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to fetch targets")
		return
	}

	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    targets,
	})
}

func (h *TargetHandler) CreateTarget(w http.ResponseWriter, r *http.Request) {
	var req CreateTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	target, err := h.queries.CreateTarget(r.Context(), db.CreateTargetParams{
		Name:            req.Name,
		Url:             req.URL,
		Method:          req.Method,
		IntervalSeconds: req.IntervalSeconds,
		TimeoutSeconds:  req.TimeoutSeconds,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to create target")
		return
	}

	WriteJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    target,
	})
}

func (h *TargetHandler) DeleteTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		WriteError(w, http.StatusBadRequest, "Target ID is required")
		return
	}

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid target ID format")
		return
	}

	if _, err := h.queries.GetTargetByID(r.Context(), id); err != nil {
		WriteError(w, http.StatusNotFound, "Target not found")
		return
	}

	if err := h.queries.SoftDeleteTarget(r.Context(), id); err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to delete target")
		return
	}

	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    "Target deleted successfully",
	})
}

type UpdateTargetRequest struct {
	CreateTargetRequest
}

func (h *TargetHandler) UpdateTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		WriteError(w, http.StatusBadRequest, "Target ID is required")
		return
	}

	var id pgtype.UUID
	if err := id.Scan(idStr); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid target ID format")
		return
	}

	var req UpdateTargetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.CreateTargetRequest.Validate(); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	target, err := h.queries.UpdateTarget(r.Context(), db.UpdateTargetParams{
		ID:              id,
		Name:            req.Name,
		Url:             req.URL,
		Method:          req.Method,
		IntervalSeconds: req.IntervalSeconds,
		TimeoutSeconds:  req.TimeoutSeconds,
	})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to update target")
		return
	}

	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    target,
	})
}
