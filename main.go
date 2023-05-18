package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

// main is the entry point for the application.
func main() {
	// Create a new MQTT client.
	buildClient()

	// Subscribe to the topic.
	// The callback function is called when a message is received.
	// The message is added to the messageCollection.
	client.Subscribe("topic/test", 0, func(client mqtt.Client, msg mqtt.Message) {
		messageCollection = append(messageCollection, NewMessage(string(msg.Payload())))
	})

	// Set up the HTTP server.
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// Define the routes for the application.
	server.GET("/messages", GetMessages)
	server.POST("/messages", AddMessage)

	// Start serving the application.
	err := server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// Message is a struct to hold a message.
type Message struct {
	ID      int       `json:"id"`
	Message string    `json:"message" binding:"required"`
	Time    time.Time `json:"timestamp"`
}

// NewMessage is the constructor for the Message struct.
func NewMessage(message string) Message {
	return Message{
		ID:      len(messageCollection) + 1,
		Message: message,
		Time:    time.Now(),
	}
}

// messageCollection is a global variable to hold the messages.
// This is a placeholder for a database.
var messageCollection []Message

// client is a global variable to hold the MQTT client.
var client mqtt.Client

// AddMessage adds a message to the messageCollection and publishes it to the topic.
// Returns a 200 status code on success.
// Returns a 400 status code on failure.
func AddMessage(c *gin.Context) {
	// instantiate a new message
	m := NewMessage("")

	// bind the JSON to the message struct
	// if there is an error, return a 400 status code
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// publish the message to the topic
	// if there is an error, return a 400 status code
	if err := publish(client, m.Message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// otherwise, return a 200 status code
	c.JSON(http.StatusOK, gin.H{
		"message": "Message published",
	})
}

// GetMessages returns a list of all collected messages.
// Returns a 200 status code.
func GetMessages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": messageCollection,
	})
}

// buildClient creates a new MQTT client.
func buildClient() {
	// broker string
	var broker = "ee58e6440f874431835beb51cc1fbd50.s2.eu.hivemq.cloud"

	// port number
	var port = 8883

	// define the options for the client
	opts := mqtt.NewClientOptions()

	// set the connection options
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetClientID("<client_name>") // set a name as you desire
	opts.SetUsername("username")      // these are the credentials that you declare for your cluster
	opts.SetPassword("Password123")

	// set callback handlers that get called on certain events
	//opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// create a new client using the options
	client = mqtt.NewClient(opts)

	// handle possible errors on connection
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// upon connection to the client, this is called
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to client")
}

// this is called when the connection to the client is lost, it prints "Connection lost" and the corresponding error
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

// publish publishes a message to the topic.
func publish(client mqtt.Client, text string) error {
	if token := client.Publish("topic/test", 0, false, text); token.Error() != nil {
		return token.Error()
	}
	return nil
}
