package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", HealthHandler)
	r.GET("/", RootHandler)
	r.POST("/data", DataHandler)

	r.Run(":8080")
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func RootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

func DataHandler(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, body)
}
