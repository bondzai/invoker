package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
		{ID: 2, Type: CronTask, CronExpr: "*/10 * * * *"}, // Run every 10 seconds
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for _, task := range tasks {
		wg.Add(1)
		go startTask(ctx, task, &wg)
	}

	// Handle interrupt signal to gracefully stop the program
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalCh
		fmt.Println("Received interrupt signal. Stopping the program...")
		cancel()
	}()

	wg.Wait()
}

func startTask(ctx context.Context, task Task, wg *sync.WaitGroup) {
	defer wg.Done()

	switch task.Type {
	case IntervalTask:
		startIntervalTask(ctx, task)
	case CronTask:
		startCronTask(ctx, task)
	}
}

func startIntervalTask(ctx context.Context, task Task) {
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

func startCronTask(ctx context.Context, task Task) {
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
