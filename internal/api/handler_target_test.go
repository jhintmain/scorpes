package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/hooneun/scorpes/internal/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestCreateTargetRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTargetRequest
		wantErr string
	}{
		{
			name: "valid request",
			req: CreateTargetRequest{
				Name:            "Test Target",
				URL:             "https://example.com",
				Method:          "GET",
				IntervalSeconds: 60,
				TimeoutSeconds:  10,
			},
			wantErr: "",
		},
		{
			name: "empty name",
			req: CreateTargetRequest{
				Name:            "",
				URL:             "https://example.com",
				Method:          "GET",
				IntervalSeconds: 60,
			},
			wantErr: "name is required",
		},
		{
			name: "whitespace only name",
			req: CreateTargetRequest{
				Name:            "   ",
				URL:             "https://example.com",
				Method:          "GET",
				IntervalSeconds: 60,
			},
			wantErr: "name is required",
		},
		{
			name: "empty url",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "",
				Method:          "GET",
				IntervalSeconds: 60,
			},
			wantErr: "url is required",
		},
		{
			name: "invalid url format",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "not-a-valid-url",
				Method:          "GET",
				IntervalSeconds: 60,
			},
			wantErr: "invalid URL format",
		},
		{
			name: "invalid http method",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "INVALID",
				IntervalSeconds: 60,
			},
			wantErr: "invalid HTTP method",
		},
		{
			name: "empty method defaults to GET",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "",
				IntervalSeconds: 60,
				TimeoutSeconds:  10,
			},
			wantErr: "",
		},
		{
			name: "lowercase method normalized to uppercase",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "post",
				IntervalSeconds: 60,
				TimeoutSeconds:  10,
			},
			wantErr: "",
		},
		{
			name: "interval seconds less than 60",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "GET",
				IntervalSeconds: 30,
			},
			wantErr: "interval seconds must be greater than or equal to 60",
		},
		{
			name: "all valid methods - POST",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "POST",
				IntervalSeconds: 60,
			},
			wantErr: "",
		},
		{
			name: "all valid methods - HEAD",
			req: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "HEAD",
				IntervalSeconds: 60,
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Validate() expected error %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("Validate() error = %q, want %q", err.Error(), tt.wantErr)
				}
			}
		})
	}
}

func TestCreateTargetRequest_Validate_DefaultValues(t *testing.T) {
	req := CreateTargetRequest{
		Name:            "Test",
		URL:             "https://example.com",
		Method:          "",
		IntervalSeconds: 60,
		TimeoutSeconds:  0,
	}

	err := req.Validate()
	if err != nil {
		t.Fatalf("Validate() unexpected error = %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("Method = %q, want %q", req.Method, "GET")
	}

	if req.TimeoutSeconds != 10 {
		t.Errorf("TimeoutSeconds = %d, want %d", req.TimeoutSeconds, 10)
	}
}

// MockQuerier implements db.Querier for testing
type MockQuerier struct {
	ListTargetsFunc      func(ctx context.Context) ([]db.Target, error)
	CreateTargetFunc     func(ctx context.Context, arg db.CreateTargetParams) (db.Target, error)
	GetTargetByIDFunc    func(ctx context.Context, id pgtype.UUID) (db.Target, error)
	SoftDeleteTargetFunc func(ctx context.Context, id pgtype.UUID) error
	UpdateTargetFunc     func(ctx context.Context, arg db.UpdateTargetParams) (db.Target, error)
}

func (m *MockQuerier) ListTargets(ctx context.Context) ([]db.Target, error) {
	if m.ListTargetsFunc != nil {
		return m.ListTargetsFunc(ctx)
	}
	return []db.Target{}, nil
}

func (m *MockQuerier) CreateTarget(ctx context.Context, arg db.CreateTargetParams) (db.Target, error) {
	if m.CreateTargetFunc != nil {
		return m.CreateTargetFunc(ctx, arg)
	}
	return db.Target{}, nil
}

func (m *MockQuerier) GetTargetByID(ctx context.Context, id pgtype.UUID) (db.Target, error) {
	if m.GetTargetByIDFunc != nil {
		return m.GetTargetByIDFunc(ctx, id)
	}
	return db.Target{}, nil
}

func (m *MockQuerier) SoftDeleteTarget(ctx context.Context, id pgtype.UUID) error {
	if m.SoftDeleteTargetFunc != nil {
		return m.SoftDeleteTargetFunc(ctx, id)
	}
	return nil
}

func (m *MockQuerier) UpdateTarget(ctx context.Context, arg db.UpdateTargetParams) (db.Target, error) {
	if m.UpdateTargetFunc != nil {
		return m.UpdateTargetFunc(ctx, arg)
	}
	return db.Target{}, nil
}

func (m *MockQuerier) GetUptimeSummary(ctx context.Context, hours int32) ([]db.GetUptimeSummaryRow, error) {
	return nil, nil
}

func (m *MockQuerier) InsertCheckResult(ctx context.Context, arg db.InsertCheckResultParams) (db.CheckResult, error) {
	return db.CheckResult{}, nil
}

func (m *MockQuerier) Ping(ctx context.Context) (int32, error) {
	return 1, nil
}

func TestTargetHandler_ListTargets(t *testing.T) {
	tests := []struct {
		name           string
		mockTargets    []db.Target
		mockErr        error
		wantStatusCode int
		wantSuccess    bool
	}{
		{
			name: "success with targets",
			mockTargets: []db.Target{
				{Name: "Target 1", Url: "https://example1.com"},
				{Name: "Target 2", Url: "https://example2.com"},
			},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name:           "success with empty list",
			mockTargets:    []db.Target{},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name:           "database error",
			mockTargets:    nil,
			mockErr:        context.DeadlineExceeded,
			wantStatusCode: http.StatusInternalServerError,
			wantSuccess:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockQuerier{
				ListTargetsFunc: func(ctx context.Context) ([]db.Target, error) {
					return tt.mockTargets, tt.mockErr
				},
			}

			handler := NewTargetHandlerWithQuerier(mock)

			req := httptest.NewRequest(http.MethodGet, "/api/targets", nil)
			rec := httptest.NewRecorder()

			handler.ListTargets(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			var resp Response
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Success != tt.wantSuccess {
				t.Errorf("success = %v, want %v", resp.Success, tt.wantSuccess)
			}
		})
	}
}

func TestTargetHandler_CreateTarget(t *testing.T) {
	tests := []struct {
		name           string
		body           any
		mockTarget     db.Target
		mockErr        error
		wantStatusCode int
		wantSuccess    bool
	}{
		{
			name: "success",
			body: CreateTargetRequest{
				Name:            "New Target",
				URL:             "https://example.com",
				Method:          "GET",
				IntervalSeconds: 60,
				TimeoutSeconds:  10,
			},
			mockTarget:     db.Target{Name: "New Target", Url: "https://example.com"},
			mockErr:        nil,
			wantStatusCode: http.StatusCreated,
			wantSuccess:    true,
		},
		{
			name:           "invalid json",
			body:           "invalid json",
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name: "validation error - empty name",
			body: CreateTargetRequest{
				Name:            "",
				URL:             "https://example.com",
				IntervalSeconds: 60,
			},
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name: "database error",
			body: CreateTargetRequest{
				Name:            "Test",
				URL:             "https://example.com",
				Method:          "GET",
				IntervalSeconds: 60,
			},
			mockErr:        context.DeadlineExceeded,
			wantStatusCode: http.StatusInternalServerError,
			wantSuccess:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockQuerier{
				CreateTargetFunc: func(ctx context.Context, arg db.CreateTargetParams) (db.Target, error) {
					return tt.mockTarget, tt.mockErr
				},
			}

			handler := NewTargetHandlerWithQuerier(mock)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/targets", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateTarget(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			var resp Response
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Success != tt.wantSuccess {
				t.Errorf("success = %v, want %v", resp.Success, tt.wantSuccess)
			}
		})
	}
}

func TestTargetHandler_DeleteTarget(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		pathValue      string
		getByIDErr     error
		deleteErr      error
		wantStatusCode int
		wantSuccess    bool
	}{
		{
			name:           "success",
			pathValue:      validUUID,
			getByIDErr:     nil,
			deleteErr:      nil,
			wantStatusCode: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name:           "empty id",
			pathValue:      "",
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:           "invalid uuid format",
			pathValue:      "invalid-uuid",
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:           "target not found",
			pathValue:      validUUID,
			getByIDErr:     context.DeadlineExceeded,
			wantStatusCode: http.StatusNotFound,
			wantSuccess:    false,
		},
		{
			name:           "delete error",
			pathValue:      validUUID,
			getByIDErr:     nil,
			deleteErr:      context.DeadlineExceeded,
			wantStatusCode: http.StatusInternalServerError,
			wantSuccess:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockQuerier{
				GetTargetByIDFunc: func(ctx context.Context, id pgtype.UUID) (db.Target, error) {
					return db.Target{}, tt.getByIDErr
				},
				SoftDeleteTargetFunc: func(ctx context.Context, id pgtype.UUID) error {
					return tt.deleteErr
				},
			}

			handler := NewTargetHandlerWithQuerier(mock)

			req := httptest.NewRequest(http.MethodDelete, "/api/targets/"+tt.pathValue, nil)
			req.SetPathValue("id", tt.pathValue)
			rec := httptest.NewRecorder()

			handler.DeleteTarget(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			var resp Response
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Success != tt.wantSuccess {
				t.Errorf("success = %v, want %v", resp.Success, tt.wantSuccess)
			}
		})
	}
}

func TestTargetHandler_UpdateTarget(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		pathValue      string
		body           any
		mockTarget     db.Target
		mockErr        error
		wantStatusCode int
		wantSuccess    bool
	}{
		{
			name:      "success",
			pathValue: validUUID,
			body: UpdateTargetRequest{
				CreateTargetRequest: CreateTargetRequest{
					Name:            "Updated Target",
					URL:             "https://updated.example.com",
					Method:          "POST",
					IntervalSeconds: 120,
					TimeoutSeconds:  15,
				},
			},
			mockTarget:     db.Target{Name: "Updated Target", Url: "https://updated.example.com"},
			mockErr:        nil,
			wantStatusCode: http.StatusOK,
			wantSuccess:    true,
		},
		{
			name:           "empty id",
			pathValue:      "",
			body:           UpdateTargetRequest{},
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:           "invalid uuid format",
			pathValue:      "invalid-uuid",
			body:           UpdateTargetRequest{},
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:           "invalid json",
			pathValue:      validUUID,
			body:           "invalid json",
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:      "validation error - empty name",
			pathValue: validUUID,
			body: UpdateTargetRequest{
				CreateTargetRequest: CreateTargetRequest{
					Name:            "",
					URL:             "https://example.com",
					IntervalSeconds: 60,
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:      "validation error - invalid url",
			pathValue: validUUID,
			body: UpdateTargetRequest{
				CreateTargetRequest: CreateTargetRequest{
					Name:            "Test",
					URL:             "not-a-valid-url",
					IntervalSeconds: 60,
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:      "validation error - interval too short",
			pathValue: validUUID,
			body: UpdateTargetRequest{
				CreateTargetRequest: CreateTargetRequest{
					Name:            "Test",
					URL:             "https://example.com",
					Method:          "GET",
					IntervalSeconds: 30,
				},
			},
			wantStatusCode: http.StatusBadRequest,
			wantSuccess:    false,
		},
		{
			name:      "database error",
			pathValue: validUUID,
			body: UpdateTargetRequest{
				CreateTargetRequest: CreateTargetRequest{
					Name:            "Test",
					URL:             "https://example.com",
					Method:          "GET",
					IntervalSeconds: 60,
				},
			},
			mockErr:        context.DeadlineExceeded,
			wantStatusCode: http.StatusInternalServerError,
			wantSuccess:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockQuerier{
				UpdateTargetFunc: func(ctx context.Context, arg db.UpdateTargetParams) (db.Target, error) {
					return tt.mockTarget, tt.mockErr
				},
			}

			handler := NewTargetHandlerWithQuerier(mock)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPut, "/api/targets/"+tt.pathValue, bytes.NewReader(bodyBytes))
			req.SetPathValue("id", tt.pathValue)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.UpdateTarget(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("status code = %d, want %d", rec.Code, tt.wantStatusCode)
			}

			var resp Response
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Success != tt.wantSuccess {
				t.Errorf("success = %v, want %v", resp.Success, tt.wantSuccess)
			}
		})
	}
}
