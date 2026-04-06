package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/version", VersionHandler)

	return r
}

func TestVersionPositive(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"returns 200"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRouter()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/version", nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if resp["version"] != "1.0.0" {
				t.Errorf("expected version '1.0.0', got '%v'", resp["version"])
			}
		})
	}
}

func TestVersionPositiveResponseStructure(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/version", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if _, ok := resp["version"]; !ok {
		t.Error("response should contain 'version' key")
	}
}

func TestVersionNegativeWrongMethod(t *testing.T) {
	methods := []string{"POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			r := setupRouter()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, "/version", nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("method %s: expected status %d, got %d", method, http.StatusNotFound, w.Code)
			}
		})
	}
}

func TestVersionNegativeNotFound(t *testing.T) {
	paths := []string{
		"/versions",
		"/ver",
		"/api/version",
		"/nonexistent",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			r := setupRouter()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", path, nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("path %s: expected status %d, got %d", path, http.StatusNotFound, w.Code)
			}
		})
	}
}
