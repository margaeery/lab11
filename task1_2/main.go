package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "lab11-task1_2/docs"
)

// @title Go API
// @version 1.0
// @description REST API with Gin and Swagger
// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", HealthHandler)
	r.POST("/echo", EchoHandler)
	r.POST("/user", UserHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
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

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler godoc
// @Summary Health check
// @Description Returns service health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{Status: "ok"})
}

// EchoHandler godoc
// @Summary Echo request body
// @Description Returns the received JSON body
// @Tags echo
// @Accept json
// @Produce json
// @Param body body object true "Any JSON body"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Router /echo [post]
func EchoHandler(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, body)
}

// UserInput represents the user creation request
type UserInput struct {
	Name string `json:"name" binding:"required" example:"Alice"`
	Age  int    `json:"age" binding:"required" example:"25"`
}

// UserResponse represents the user creation response
type UserResponse struct {
	Message string    `json:"message"`
	User    UserInput `json:"user"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// UserHandler godoc
// @Summary Create a user
// @Description Creates a user with name and age validation
// @Tags user
// @Accept json
// @Produce json
// @Param input body UserInput true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /user [post]
func UserHandler(c *gin.Context) {
	var input UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Age < 18 || input.Age > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "age must be between 18 and 100"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": input})
}
