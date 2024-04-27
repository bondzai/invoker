package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// Consumer interface defines the behavior of a RabbitMQ consumer
type Consumer interface {
	Consume() error
}

// RabbitMQConsumer struct holds the necessary fields for consuming messages from RabbitMQ
type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

// NewRabbitMQConsumer creates a new RabbitMQ consumer instance
func NewRabbitMQConsumer(url, queueName string) (*RabbitMQConsumer, error) {
	consumer := &RabbitMQConsumer{}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	consumer.connection = conn

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	consumer.channel = ch

	q, err := ch.QueueDeclare(
		queueName, // Queue name
		true,      // Durable
		false,     // Delete when unused
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return nil, err
	}
	consumer.queue = q

	return consumer, nil
}

// Consume starts consuming messages from RabbitMQ queue
func (c *RabbitMQConsumer) Consume() error {
	msgs, err := c.channel.Consume(
		c.queue.Name, // Queue
		"",           // Consumer
		true,         // Auto Ack
		false,        // Exclusive
		false,        // No local
		false,        // No wait
		nil,          // Args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data interface{}
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("Error decoding message body: %v", err)
				continue
			}
			// Process the received message, for now, just print it
			log.Printf("Received a message: %v", data)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
