package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/howeyc/gopass"
)

// messageCollection is a global variable to hold the messages.
// This is a placeholder for a database.
var messageCollection []Message

// client is a global variable to hold the MQTT client.
var client mqtt.Client

// env is a global variable to hold the environment variables.
var env Env

// init is called before the application starts.
// It prompts the user for the environment variables.
func init() {
	fmt.Println("Starting the application...")
	fmt.Println("Enter the url of your broker(e.g. broker.hivemq.com):")
	if _, err := fmt.Scanln(&env.broker); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enter the port of your broker(e.g. 1883):")
	if _, err := fmt.Scanln(&env.port); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enter the topic you want to subscribe to(e.g. topic/test):")
	if _, err := fmt.Scanln(&env.topic); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enter the username of your broker:")
	if _, err := fmt.Scanln(&env.username); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Enter the password of your broker:")
	password, err := gopass.GetPasswd()
	if err != nil {
		log.Fatal(err)
	}
	env.password = string(password)
}

// main is the entry point for the application.
func main() {
	// Create a new MQTT client.
	buildClient()

	// Subscribe to the topic.
	// The callback function is called when a Message is received.
	// The Message is added to the messageCollection.
	client.Subscribe(env.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		messageCollection = append(messageCollection, newMessage(string(msg.Payload())))
	})

	// Set up the HTTP server.
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	// Define the routes for the application.
	server.GET("/messages", getMessages)
	server.POST("/messages", addMessage)

	// Start serving the application.
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

// Message is a struct to hold a message.
type Message struct {
	ID      int       `json:"id"`
	Message string    `json:"message" binding:"required"`
	Time    time.Time `json:"timestamp"`
}

// newMessage is the constructor for the Message struct.
func newMessage(message string) Message {
	return Message{
		ID:      len(messageCollection) + 1,
		Message: message,
		Time:    time.Now(),
	}
}

// Env is a struct to hold the environment variables.
type Env struct {
	broker   string
	port     int
	topic    string
	username string
	password string
}

// addMessage adds a Message to the messageCollection and publishes it to the topic.
// Returns a 200 status code on success.
// Returns a 400 status code on failure.
func addMessage(c *gin.Context) {
	// instantiate a new Message
	m := newMessage("")

	// bind the JSON to the Message struct
	// if there is an error, return a 400 status code
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// publish the Message to the topic
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

// getMessages returns a list of all collected messages.
// Returns a 200 status code.
func getMessages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": messageCollection,
	})
}

// buildClient creates a new MQTT client.
func buildClient() {
	// instantiate a new ClientOptions struct
	opts := mqtt.NewClientOptions()

	// set the connection options
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", env.broker, env.port))
	opts.SetUsername(env.username)
	opts.SetPassword(env.password)

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
	fmt.Println("Successfully connected")
	fmt.Println("Broker:", env.broker)
	fmt.Println("Subscribing to topic:", env.topic)
	fmt.Println("http server listening on port:", env.port)
	fmt.Println("===================================")
	fmt.Println("to publish a message, send a POST request to /messages")
	fmt.Println("the request body should be a JSON object with the following structure:")
	fmt.Println("{")
	fmt.Println("  \"message\": \"your message\"")
	fmt.Println("}")
	fmt.Println("===================================")
	fmt.Println("to get all messages, send a GET request to /messages")
	fmt.Println("===================================")
	fmt.Println("to exit the application, press CTRL+C")
	fmt.Println("===================================")
	fmt.Println("HTTP Request Log:")
}

// this is called when the connection to the client is lost, it prints "Connection lost" and the corresponding error
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

// publish publishes a Message to the topic.
func publish(client mqtt.Client, text string) error {
	if token := client.Publish("topic/test", 0, false, text); token.Error() != nil {
		return token.Error()
	}
	return nil
}
