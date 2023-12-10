package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/bondzai/invoker/internal/mock"
	"github.com/bondzai/invoker/internal/task"
)

const (
	httpPort  = 8080
	httpRoute = "/ping"
)

var (
	numTasks    = 20000
	numTasksMux sync.Mutex
)

func main() {
	tasks := mock.GenerateTasks(numTasks)

	var wg sync.WaitGroup

	taskManagers := map[task.TaskType]task.TaskManager{
		task.IntervalTask: &task.IntervalTaskManager{},
		task.CronTask:     &task.CronTaskManager{},
	}

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
		http.HandleFunc(httpRoute, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Invoker is running...")
		})

		http.HandleFunc("/setNumTasks", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var requestBody struct {
				NumTasks int `json:"numTasks"`
			}

			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
				return
			}

			numTasksMux.Lock()
			numTasks = requestBody.NumTasks
			numTasksMux.Unlock()

			newTasks := mock.GenerateTasks(numTasks)
			for _, t := range newTasks {
				wg.Add(1)
				go func(task task.Task) {
					defer wg.Done()
					taskManagers[task.Type].Start(ctx, task, &wg)
				}(t)
			}

			fmt.Fprintf(w, "numTasks updated to %d", requestBody.NumTasks)
		})

		http.HandleFunc("/getNumTasks", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			numTasksMux.Lock()
			currentNumTasks := numTasks
			numTasksMux.Unlock()

			fmt.Fprintf(w, "Current number of tasks: %d", currentNumTasks)
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
		if err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
			cancel()
		}
	}()

	// Start tasks invoke loop
	for _, t := range tasks {
		wg.Add(1)
		go func(task task.Task) {
			defer wg.Done()
			taskManagers[task.Type].Start(ctx, task, &wg)
		}(t)
	}

	wg.Wait()
}
