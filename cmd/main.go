package main

import (
	"context"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/shutdown"
	"github.com/bondzai/invoker/internal/signalhandler"
	"github.com/bondzai/invoker/internal/task"
)

func main() {
	tasks := []task.Task{
		{ID: 1, Type: task.IntervalTask, Interval: 5 * time.Second},
		{ID: 2, Type: task.CronTask, CronExpr: "*/10 * * * *"},
	}

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
