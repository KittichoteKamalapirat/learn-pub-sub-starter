package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/KittichoteKamalapirat/learn-pub-sub-starter/internal/pubsub"
	"github.com/KittichoteKamalapirat/learn-pub-sub-starter/internal/routing"

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

	// Create a PlayingState message
	state := routing.PlayingState{IsPaused: true}
	body, err := json.Marshal(state)
	if err != nil {
		fmt.Println("Failed to marshal JSON:", err)
		return
	}

	// Publish the message
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, "pause", body)
	if err != nil {
		fmt.Println("Failed to publish message:", err)
		return
	}

	// Wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
}
