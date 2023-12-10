package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/bondzai/invoker/internal/gracefulshutdown"
	"github.com/bondzai/invoker/internal/mock"
	"github.com/bondzai/invoker/internal/task"
)

const (
	numTasks  = 10000
	httpPort  = 8080
	httpRoute = "/"
)

func main() {
	// Generate mock tasks
	tasks := mock.GenerateTasks(numTasks)

	// Create a graceful shutdown manager
	shutdownManager := gracefulshutdown.NewManager()
	shutdownManager.StartSignalHandling()

	var wg sync.WaitGroup

	// Task managers for different types of tasks
	taskManagers := map[task.TaskType]task.TaskManager{
		task.IntervalTask: &task.IntervalTaskManager{},
		task.CronTask:     &task.CronTaskManager{},
	}

	// Start HTTP server in a goroutine
	go func() {
		http.HandleFunc(httpRoute, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "HTTP server is running")
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
		if err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			shutdownManager.Shutdown() // Trigger shutdown on HTTP server error
		}
	}()

	// Start tasks
	for _, t := range tasks {
		wg.Add(1)
		go func(task task.Task) {
			defer wg.Done()
			taskManagers[task.Type].Start(shutdownManager.Context(), task, shutdownManager.WaitGroup(), shutdownManager)
		}(t)
	}

	// Wait for tasks to finish
	wg.Wait()

	// Shutdown HTTP server
	shutdownManager.Shutdown()
}
