package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/bondzai/invoker/internal/mock"
	"github.com/bondzai/invoker/internal/task"
)

const (
	numTasks  = 10000
	httpPort  = 8080
	httpRoute = "/ping"
)

func main() {
	tasks := mock.GenerateTasks(numTasks)

	var wg sync.WaitGroup

	taskManagers := map[task.TaskType]task.TaskManager{
		task.IntervalTask: &task.IntervalTaskManager{},
		task.CronTask:     &task.CronTaskManager{},
	}

	// Start HTTP server in a goroutine
	go func() {
		http.HandleFunc(httpRoute, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Invoker is running...")
		})

		err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
		if err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
		}
	}()

	// Start tasks invoke loop
	for _, t := range tasks {
		wg.Add(1)
		go func(task task.Task) {
			defer wg.Done()
			taskManagers[task.Type].Start(context.Background(), task, &wg)
		}(t)
	}

	wg.Wait()
}
