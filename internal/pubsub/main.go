package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON publishes a JSON message to a RabbitMQ exchange.
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	fmt.Println("Publishing JSON message")
	body, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println("Publishing JSON message 2", body)

	return ch.PublishWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("failed to open a channel: %w", err)
	}

	queue, err := ch.QueueDeclare(queueName, simpleQueueType == 0, simpleQueueType == 1, simpleQueueType == 1, false, nil)

	if err != nil {
		fmt.Println("Error creating queue:", err)
	}
	ch.QueueBind(queueName, key, exchange, simpleQueueType == 1, nil)

	return ch, queue, nil
}
