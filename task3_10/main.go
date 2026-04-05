package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port         string
	Env          string
	AppName      string
	MaxBodySize  int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "go-config-service"
	}

	maxBodySize, err := strconv.Atoi(os.Getenv("MAX_BODY_SIZE"))
	if err != nil || maxBodySize <= 0 {
		maxBodySize = 10
	}

	readTimeout, err := time.ParseDuration(os.Getenv("READ_TIMEOUT"))
	if err != nil {
		readTimeout = 5 * time.Second
	}

	writeTimeout, err := time.ParseDuration(os.Getenv("WRITE_TIMEOUT"))
	if err != nil {
		writeTimeout = 10 * time.Second
	}

	return Config{
		Port:         port,
		Env:          env,
		AppName:      appName,
		MaxBodySize:  maxBodySize,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func main() {
	cfg := LoadConfig()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", HealthHandler)
	r.GET("/config", func(c *gin.Context) {
		ConfigHandler(c, cfg)
	})
	r.POST("/echo", EchoHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		log.Printf("[%s] %s listening on :%s", cfg.Env, cfg.AppName, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}

	log.Println("server stopped")
}

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

type ConfigResponse struct {
	AppName      string `json:"app_name"`
	Env          string `json:"env"`
	Port         string `json:"port"`
	MaxBodySize  int    `json:"max_body_size_mb"`
	ReadTimeout  string `json:"read_timeout"`
	WriteTimeout string `json:"write_timeout"`
}

func ConfigHandler(c *gin.Context, cfg Config) {
	c.JSON(http.StatusOK, ConfigResponse{
		AppName:      cfg.AppName,
		Env:          cfg.Env,
		Port:         cfg.Port,
		MaxBodySize:  cfg.MaxBodySize,
		ReadTimeout:  cfg.ReadTimeout.String(),
		WriteTimeout: cfg.WriteTimeout.String(),
	})
}

type EchoRequest struct {
	Message string `json:"message"`
}

type EchoResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Server    string `json:"server"`
}

func EchoHandler(c *gin.Context) {
	var req EchoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "go-config-service"
	}

	c.JSON(http.StatusOK, EchoResponse{
		Message:   req.Message,
		Timestamp: time.Now().Format(time.RFC3339),
		Server:    appName,
	})
}

func GetDefaultConfig() Config {
	return Config{
		Port:         "8080",
		Env:          "development",
		AppName:      "go-config-service",
		MaxBodySize:  10,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
