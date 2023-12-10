package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/shutdown"
	"github.com/bondzai/invoker/internal/signalhandler"
	"github.com/robfig/cron/v3"
)

type TaskType int

const (
	IntervalTask TaskType = iota
	CronTask
)

type Task struct {
	ID       int
	Type     TaskType
	Interval time.Duration
	CronExpr string
}

func main() {
	tasks := []Task{
		{ID: 1, Type: IntervalTask, Interval: 5 * time.Second},
		{ID: 2, Type: CronTask, CronExpr: "*/10 * * * *"},
	}

	shutdownManager := shutdown.NewGracefulShutdownManager()
	signalHandler := signalhandler.NewSignalHandler()
	signalHandler.Start()

	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go startTask(context.Background(), task, &wg, shutdownManager)
	}

	wg.Wait()
	shutdownManager.Shutdown()
}
func startTask(ctx context.Context, task Task, wg *sync.WaitGroup, shutdownManager shutdown.ShutdownManager) {
	defer wg.Done()

	switch task.Type {
	case IntervalTask:
		startIntervalTask(ctx, task, shutdownManager)
	case CronTask:
		startCronTask(ctx, task, shutdownManager)
	}
}

func startIntervalTask(ctx context.Context, task Task, shutdownManager shutdown.ShutdownManager) {
	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("Interval Task %d: Triggered at %v\n", task.ID, time.Now())
			// Add your interval task-specific logic here

		case <-ctx.Done():
			fmt.Printf("Interval Task %d: Stopping...\n", task.ID)
			return
		}
	}
}

func startCronTask(ctx context.Context, task Task, shutdownManager shutdown.ShutdownManager) {
	c := cron.New()
	defer c.Stop()

	_, err := c.AddFunc(task.CronExpr, func() {
		fmt.Printf("Cron Task %d: Triggered at %v\n", task.ID, time.Now())
		// Add your cron task-specific logic here
	})
	if err != nil {
		fmt.Printf("Cron Task %d: Error adding cron expression: %v\n", task.ID, err)
		return
	}

	c.Start()

	<-ctx.Done()
	fmt.Printf("Cron Task %d: Stopping...\n", task.ID)
}
