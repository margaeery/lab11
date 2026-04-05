package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/health", HealthHandler)
	r.POST("/echo", EchoHandler)
	r.POST("/user", UserHandler)
	return r
}

func TestHealthHandler(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestEchoHandler(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid json object",
			body:       `{"hello":"world"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty object",
			body:       `{}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "nested object",
			body:       `{"user":{"name":"Alice","age":25}}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid json",
			body:       `{invalid}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not json",
			body:       `hello`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d, body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Logf("response is not a JSON object (may be array): %v", err)
				}
			}
		})
	}
}

func TestUserHandler(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid user",
			body:       `{"name":"Alice","age":25}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "valid user boundary 18",
			body:       `{"name":"Bob","age":18}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "valid user boundary 100",
			body:       `{"name":"Carol","age":100}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "age below minimum",
			body:       `{"name":"Dave","age":17}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "age above maximum",
			body:       `{"name":"Eve","age":101}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing name",
			body:       `{"age":25}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing age",
			body:       `{"name":"Frank"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid json",
			body:       `{bad}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty name string",
			body:       `{"name":"","age":30}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative age",
			body:       `{"name":"Grace","age":-5}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "zero age",
			body:       `{"name":"Henry","age":0}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d, body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusCreated {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if resp["message"] != "user created" {
					t.Errorf("expected message 'user created', got '%v'", resp["message"])
				}
			}
		})
	}
}
