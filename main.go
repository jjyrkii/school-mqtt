package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Message is a struct to hold a message.
type Message struct {
	Message string    `json:"message" binding:"required"`
	Time    time.Time `json:"timestamp"`
}

func NewMessage(message string) Message {
	return Message{
		Message: message,
		Time:    time.Now(),
	}
}

// MessageCollection is a struct to hold a collection of messages.
type MessageCollection struct {
	Messages []Message `json:"messages"`
}

// collection is a global variable to hold the messages.
// This is a placeholder for a database.
var collection MessageCollection

// main is the entry point for the application.
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
	m := NewMessage("pong")
	c.JSON(http.StatusOK, gin.H{
		"body": m,
	})
}

// AddMessage adds a message to the collection.
// Returns a 200 status code if successful.
// Returns a 400 status code if the message is missing or not a string.
func AddMessage(c *gin.Context) {
	m := NewMessage("")
	err := c.ShouldBindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Message is required and must be a string.",
		})
		return
	}

	collection.Messages = append(collection.Messages, m)

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
