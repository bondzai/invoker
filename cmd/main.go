package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/bondzai/invoker/internal/mock"
	"github.com/bondzai/invoker/internal/task"
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

	// Get tasks from mock package
	tasks := mock.GetTasks()

	// Map task types to task managers
	taskManagers := map[task.TaskType]task.TaskManager{
		task.IntervalTask: &task.IntervalTaskManager{},
		task.CronTask:     &task.CronTaskManager{},
	}

	// Start tasks invoke loop
	for _, t := range *tasks {
		wg.Add(1)
		go func(task task.Task) {
			defer wg.Done()
			taskManagers[task.Type].Start(ctx, task, &wg)
		}(t)
	}

	wg.Wait()
}
