package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// main is the entry point for the application.
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
