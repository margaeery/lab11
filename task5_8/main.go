package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": "1.0.0",
		})
	})

	r.Run(":8080")
}
