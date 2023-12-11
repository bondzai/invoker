package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/bondzai/invoker/internal/task"
	"github.com/streadway/amqp"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown using a goroutine
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt)

		<-sigCh

		fmt.Println("\nReceived interrupt signal. Initiating graceful shutdown...")
		cancel()

		wg.Wait()

		fmt.Println("Shutdown complete.")
		os.Exit(0)
	}()

	// Start HTTP server in a goroutine
	go func() {
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Invoker is running...")
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
		if err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			cancel()
		}
	}()

	// Map task types to task managers
	taskManagers := map[task.TaskType]task.TaskManager{
		task.IntervalTask: &task.IntervalTaskManager{},
		task.CronTask:     &task.CronTaskManager{},
	}

	// Create a RabbitMQ connection and channel
	conn, err := amqp.Dial("amqp://guest:guest@172.21.0.2:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := declareQueue(ch)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Start tasks invoke loop
	consumeTasks(ctx, ch, q.Name, taskManagers, &wg)
}

func declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"tasks",
		false,
		false,
		false,
		false,
		nil,
	)
}

func consumeTasks(ctx context.Context, ch *amqp.Channel, queueName string, taskManagers map[task.TaskType]task.TaskManager, wg *sync.WaitGroup) {
	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			var t task.Task
			err := json.Unmarshal(msg.Body, &t)
			if err != nil {
				log.Printf("Failed to unmarshal task: %v", err)
				continue
			}

			wg.Add(1)
			go func(task task.Task) {
				defer wg.Done()
				taskManagers[task.Type].Start(ctx, task, wg)
			}(t)
		}
	}
}
