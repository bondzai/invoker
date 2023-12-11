// rabbitmq/consumer.go
package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Config holds the configuration for connecting to RabbitMQ
type Config struct {
	URL      string
	Queue    string
	Exchange string
}

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	config     Config
	connection *amqp.Connection
	channel    *amqp.Channel
	messages   <-chan amqp.Delivery
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(config Config) *Consumer {
	conn, err := amqp.Dial(config.URL)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	_, err = ch.QueueDeclare(
		config.Queue, // Queue name
		true,         // Durable
		false,        // Delete when unused
		false,        // Exclusive
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		config.Queue, // Queue
		"",           // Consumer
		false,        // Auto Ack
		false,        // Exclusive
		false,        // No Local
		false,        // No Wait
		nil,          // Args
	)
	if err != nil {
		panic(err)
	}

	return &Consumer{
		config:     config,
		connection: conn,
		channel:    ch,
		messages:   msgs,
	}
}

// Messages returns the channel for consuming messages
func (c *Consumer) Messages() <-chan amqp.Delivery {
	return c.messages
}

// Close closes the RabbitMQ connection and channel
func (c *Consumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
