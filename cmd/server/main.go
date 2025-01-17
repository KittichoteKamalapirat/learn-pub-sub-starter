package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/KittichoteKamalapirat/learn-pub-sub-starter/internal/pubsub"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// RabbitMQ connection string
	rabbit_url := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbit_url)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ")
	fmt.Println("Starting Peril server...")

	// Create a new channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel:", err)
		return
	}
	defer ch.Close()

	// Declare and bind the durable queue named game_logs
	queueName := "game_logs"
	exchange := "peril_topic"   // Use the new exchange
	routingKey := "game_logs.*" // Routing key for the queue
	simpleQueueType := 0        // Set to 0 for durable

	// Call DeclareAndBind
	_, _, err = pubsub.DeclareAndBind(conn, exchange, queueName, routingKey, simpleQueueType)
	if err != nil {
		fmt.Println("Failed to declare and bind queue:", err)
		return
	}

	// Wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
}
