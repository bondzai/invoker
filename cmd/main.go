package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/mock"
	"github.com/bondzai/invoker/internal/task"
)

func init() {
	fmt.Println("Invoker is starting...")
}

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	go handleGracefulShutdown(cancel, &wg)

	server := api.NewHttpServer()
	go server.Start(ctx)

	// Get task managers mapping
	taskManagers := *task.NewTaskManagers()

	// Start tasks invoke loop
	for _, t := range *mock.GetTasks() {
		wg.Add(1)
		go func(task task.Task) {
			defer wg.Done()
			taskManagers[task.Type].Start(ctx, task, &wg, nil)
		}(t)
	}

	wg.Wait()
}

// Handle graceful shutdown using a goroutine
func handleGracefulShutdown(cancel context.CancelFunc, wg *sync.WaitGroup) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	<-sigCh

	fmt.Println("\nReceived interrupt signal. Initiating graceful shutdown...")
	cancel()

	wg.Wait()

	fmt.Println("Shutdown complete.")
	os.Exit(0)
}
