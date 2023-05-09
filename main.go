package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Message is a struct to hold a message.
type Message struct {
	Message string `json:"message"`
}

// MessageCollection is a struct to hold a collection of messages.
type MessageCollection struct {
	Messages []Message `json:"messages"`
}

// main is the entry point for the application.
func main() {

	// Server settings
	r := gin.Default()

	// Define the routes for the application.
	r.GET("/ping", GetPong)

	// Start serving the application.
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// GetPong returns a simple message.
func GetPong(c *gin.Context) {
	message := Message{"pong"}
	c.JSON(http.StatusOK, gin.H{
		"message": message.Message,
	})
}
