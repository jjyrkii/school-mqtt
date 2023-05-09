package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Message is a struct to hold a message.
type Message struct {
	Message string `json:"message" binding:"required"`
}

// MessageCollection is a struct to hold a collection of messages.
type MessageCollection struct {
	Messages []Message `json:"messages"`
}

// main is the entry point for the application.
var collection MessageCollection

func main() {

	// Server settings
	r := gin.Default()

	// Define the routes for the application.
	r.GET("/ping", GetPong)
	r.GET("/messages", GetMessages)
	r.POST("/messages", AddMessage)

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
		"body": message,
	})
}

// AddMessage adds a message to the collection.
// Returns a 200 status code if successful.
// Returns a 400 status code if the message is missing or not a string.
func AddMessage(c *gin.Context) {
	var json Message
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Message is required and must be a string.",
		})
		return
	}

	collection.Messages = append(collection.Messages, json)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message successfully added.",
	})
}

// GetMessages returns a list of messages.
func GetMessages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": collection,
	})
}
