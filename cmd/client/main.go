package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/KittichoteKamalapirat/learn-pub-sub-starter/internal/gamelogic"
	"github.com/KittichoteKamalapirat/learn-pub-sub-starter/internal/pubsub"
	"github.com/KittichoteKamalapirat/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	// RabbitMQ connection string
	rabbit_url := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbit_url)
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Error getting username:", err)
		return
	}

	// Declare and bind the queue
	queueName := fmt.Sprintf("pause.%s", username) // Create the queue name
	exchange := routing.ExchangePerilDirect        // Use the constant for the exchange
	routingKey := routing.PauseKey                 // Use the constant for the routing key
	simpleQueueType := 1                           // Set to 1 for transient

	// Call DeclareAndBind
	ch, _, err := pubsub.DeclareAndBind(conn, exchange, queueName, routingKey, simpleQueueType)
	if err != nil {
		fmt.Println("Failed to declare and bind queue:", err)
		return
	}
	defer ch.Close() // Ensure the channel is closed when done

	// Create a new game state
	gameState := gamelogic.NewGameState(username)

	// Start the REPL loop
	for {
		// Get user input
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		command := words[0]
		switch command {
		case "spawn":
			if err := gameState.CommandSpawn(words); err != nil {
				fmt.Println("Error:", err)
			}

		case "move":
			if army, err := gameState.CommandMove(words); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Move successful.")
				fmt.Println(army)
			}

		case "status":
			gameState.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			fmt.Println("Don't understand that command")
		}
	}

	// Wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
}
