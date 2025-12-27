package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")

	v1.GET("/messages", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, Friends! This is a little message from me."})
	})

	_ = r.Run(":8080")
}
