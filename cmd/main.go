package main

import (
	"context"
	"sync"

	"github.com/bondzai/invoker/internal/mock"
	"github.com/bondzai/invoker/internal/shutdown"
	"github.com/bondzai/invoker/internal/signalhandler"
	"github.com/bondzai/invoker/internal/task"
)

const numTasks = 10000

func main() {
	tasks := mock.GenerateTasks(numTasks)

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
