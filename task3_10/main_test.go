package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	cfg := LoadConfig()
	r.GET("/health", HealthHandler)
	r.GET("/config", func(c *gin.Context) {
		ConfigHandler(c, cfg)
	})
	r.POST("/echo", EchoHandler)
	return r
}

func TestHealthHandler(t *testing.T) {
	t.Run("returns 200 and ok status", func(t *testing.T) {
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
	})

	t.Run("returns 404 for unregistered method", func(t *testing.T) {
		router := setupRouter()
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestConfigHandler(t *testing.T) {
	t.Run("returns config with defaults", func(t *testing.T) {
		router := setupRouter()
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var resp ConfigResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		if resp.EnvVars["PORT"] != "8080" {
			t.Errorf("expected PORT '8080', got '%s'", resp.EnvVars["PORT"])
		}
		if resp.EnvVars["APP_ENV"] != "development" {
			t.Errorf("expected APP_ENV 'development', got '%s'", resp.EnvVars["APP_ENV"])
		}
		if resp.EnvVars["APP_NAME"] != "go-config-service" {
			t.Errorf("expected APP_NAME 'go-config-service', got '%s'", resp.EnvVars["APP_NAME"])
		}
	})

	t.Run("masks sensitive env vars", func(t *testing.T) {
		os.Setenv("DB_PASSWORD", "secret123")
		os.Setenv("API_KEY", "key-abc")
		defer func() {
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("API_KEY")
		}()

		router := setupRouter()
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp ConfigResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		if resp.EnvVars["DB_PASSWORD"] != "****" {
			t.Errorf("expected DB_PASSWORD masked, got %q", resp.EnvVars["DB_PASSWORD"])
		}
		if resp.EnvVars["API_KEY"] != "****" {
			t.Errorf("expected API_KEY masked, got %q", resp.EnvVars["API_KEY"])
		}
	})

	t.Run("shows non-sensitive env vars", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		os.Setenv("APP_NAME", "my-app")
		defer func() {
			os.Unsetenv("APP_ENV")
			os.Unsetenv("APP_NAME")
		}()

		router := setupRouter()
		req := httptest.NewRequest(http.MethodGet, "/config", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp ConfigResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		if resp.EnvVars["APP_ENV"] != "production" {
			t.Errorf("expected APP_ENV 'production', got %q", resp.EnvVars["APP_ENV"])
		}
		if resp.EnvVars["APP_NAME"] != "my-app" {
			t.Errorf("expected APP_NAME 'my-app', got %q", resp.EnvVars["APP_NAME"])
		}
	})

	t.Run("returns 404 for unregistered method", func(t *testing.T) {
		router := setupRouter()
		req := httptest.NewRequest(http.MethodPost, "/config", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestEchoHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid message",
			body:       `{"message":"hello"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty message string",
			body:       `{"message":""}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "unicode message",
			body:       `{"message":"привет мир"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid json syntax",
			body:       `{bad}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "plain text instead of json",
			body:       `hello`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "array instead of object",
			body:       `[1,2,3]`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter()
			req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d, body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.wantStatus == http.StatusOK {
				var resp EchoResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if resp.Timestamp == "" {
					t.Error("expected timestamp to be non-empty")
				}
				if resp.Server == "" {
					t.Error("expected server name to be non-empty")
				}
			}
		})
	}

	t.Run("returns 404 for unregistered method", func(t *testing.T) {
		router := setupRouter()
		req := httptest.NewRequest(http.MethodGet, "/echo", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		check       func(t *testing.T, cfg Config)
	}{
		{
			name:    "all defaults",
			envVars: map[string]string{},
			check: func(t *testing.T, cfg Config) {
				if cfg.Port != "8080" {
					t.Errorf("port: got %q, want %q", cfg.Port, "8080")
				}
				if cfg.Env != "development" {
					t.Errorf("env: got %q, want %q", cfg.Env, "development")
				}
				if cfg.AppName != "go-config-service" {
					t.Errorf("app_name: got %q, want %q", cfg.AppName, "go-config-service")
				}
				if cfg.MaxBodySize != 10 {
					t.Errorf("max_body_size: got %d, want %d", cfg.MaxBodySize, 10)
				}
				if cfg.ReadTimeout != 5*time.Second {
					t.Errorf("read_timeout: got %v, want %v", cfg.ReadTimeout, 5*time.Second)
				}
				if cfg.WriteTimeout != 10*time.Second {
					t.Errorf("write_timeout: got %v, want %v", cfg.WriteTimeout, 10*time.Second)
				}
			},
		},
		{
			name: "all custom",
			envVars: map[string]string{
				"PORT":          "9090",
				"APP_ENV":       "production",
				"APP_NAME":      "my-app",
				"MAX_BODY_SIZE": "50",
				"READ_TIMEOUT":  "3s",
				"WRITE_TIMEOUT": "15s",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.Port != "9090" {
					t.Errorf("port: got %q, want %q", cfg.Port, "9090")
				}
				if cfg.Env != "production" {
					t.Errorf("env: got %q, want %q", cfg.Env, "production")
				}
				if cfg.AppName != "my-app" {
					t.Errorf("app_name: got %q, want %q", cfg.AppName, "my-app")
				}
				if cfg.MaxBodySize != 50 {
					t.Errorf("max_body_size: got %d, want %d", cfg.MaxBodySize, 50)
				}
				if cfg.ReadTimeout != 3*time.Second {
					t.Errorf("read_timeout: got %v, want %v", cfg.ReadTimeout, 3*time.Second)
				}
				if cfg.WriteTimeout != 15*time.Second {
					t.Errorf("write_timeout: got %v, want %v", cfg.WriteTimeout, 15*time.Second)
				}
			},
		},
		{
			name: "invalid max_body_size falls back",
			envVars: map[string]string{
				"MAX_BODY_SIZE": "invalid",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.MaxBodySize != 10 {
					t.Errorf("max_body_size: got %d, want %d", cfg.MaxBodySize, 10)
				}
			},
		},
		{
			name: "invalid read_timeout falls back",
			envVars: map[string]string{
				"READ_TIMEOUT": "bad",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.ReadTimeout != 5*time.Second {
					t.Errorf("read_timeout: got %v, want %v", cfg.ReadTimeout, 5*time.Second)
				}
			},
		},
		{
			name: "invalid write_timeout falls back",
			envVars: map[string]string{
				"WRITE_TIMEOUT": "wrong",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.WriteTimeout != 10*time.Second {
					t.Errorf("write_timeout: got %v, want %v", cfg.WriteTimeout, 10*time.Second)
				}
			},
		},
		{
			name: "negative max_body_size falls back",
			envVars: map[string]string{
				"MAX_BODY_SIZE": "-5",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.MaxBodySize != 10 {
					t.Errorf("max_body_size: got %d, want %d", cfg.MaxBodySize, 10)
				}
			},
		},
		{
			name: "zero max_body_size falls back",
			envVars: map[string]string{
				"MAX_BODY_SIZE": "0",
			},
			check: func(t *testing.T, cfg Config) {
				if cfg.MaxBodySize != 10 {
					t.Errorf("max_body_size: got %d, want %d", cfg.MaxBodySize, 10)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			cfg := LoadConfig()

			for k := range tt.envVars {
				os.Unsetenv(k)
			}

			tt.check(t, cfg)
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	cfg := GetDefaultConfig()

	if cfg.Port != "8080" {
		t.Errorf("default port: got %q, want %q", cfg.Port, "8080")
	}
	if cfg.Env != "development" {
		t.Errorf("default env: got %q, want %q", cfg.Env, "development")
	}
	if cfg.AppName != "go-config-service" {
		t.Errorf("default app_name: got %q, want %q", cfg.AppName, "go-config-service")
	}
	if cfg.MaxBodySize != 10 {
		t.Errorf("default max_body_size: got %d, want %d", cfg.MaxBodySize, 10)
	}
	if cfg.ReadTimeout != 5*time.Second {
		t.Errorf("default read_timeout: got %v, want %v", cfg.ReadTimeout, 5*time.Second)
	}
	if cfg.WriteTimeout != 10*time.Second {
		t.Errorf("default write_timeout: got %v, want %v", cfg.WriteTimeout, 10*time.Second)
	}
}

func TestNotFoundRoute(t *testing.T) {
	router := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestEchoResponseContainsServerName(t *testing.T) {
	os.Setenv("APP_NAME", "test-server")
	defer os.Unsetenv("APP_NAME")

	router := setupRouter()
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(`{"message":"hi"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp EchoResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if resp.Server != "test-server" {
		t.Errorf("expected server 'test-server', got '%s'", resp.Server)
	}
}
