package main

import (
	"context"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/shutdown"
	"github.com/bondzai/invoker/internal/signalhandler"
	"github.com/bondzai/invoker/internal/task"
)

const numTasks = 10000

func main() {
	tasks := generateTasks(numTasks)

	shutdownManager := shutdown.NewGracefulShutdownManager()
	signalHandler := signalhandler.NewSignalHandler()
	signalHandler.Start()

	var wg sync.WaitGroup

	taskManagers := map[task.TaskType]task.TaskManager{
		task.IntervalTask: &task.IntervalTaskManager{},
		task.CronTask:     &task.CronTaskManager{},
	}

	for _, task := range tasks {
		wg.Add(1)
		go taskManagers[task.Type].Start(context.Background(), task, &wg, shutdownManager)
	}

	wg.Wait()
	shutdownManager.Shutdown()
}

// generateTasks generates a slice of tasks with the specified number.
// The generated tasks alternate between interval tasks and cron tasks.
// For interval tasks, the interval is set to 5 seconds.
// For cron tasks, the cron expression is set to "*/10 * * * *".
//
//	ex tasks := []task.Task{
//		{ID: 1, Type: task.IntervalTask, Interval: 5 * time.Second},
//		{ID: 2, Type: task.CronTask, CronExpr: "*/10 * * * *"},
//	}
func generateTasks(numTasks int) []task.Task {
	tasks := make([]task.Task, numTasks)

	for i := 0; i < numTasks; i++ {
		if i%2 == 0 {
			tasks[i] = task.Task{
				ID:       i + 1,
				Type:     task.IntervalTask,
				Interval: 5 * time.Second,
			}
		} else {
			tasks[i] = task.Task{
				ID:       i + 1,
				Type:     task.CronTask,
				CronExpr: "*/10 * * * *",
			}
		}
	}

	return tasks
}
