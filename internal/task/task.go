// task.go
package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/shutdown"
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

type TaskManager interface {
	Start(ctx context.Context, task Task, wg *sync.WaitGroup, shutdownManager shutdown.ShutdownManager)
}

type IntervalTaskManager struct{}

func (m *IntervalTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup, shutdownManager shutdown.ShutdownManager) {
	defer wg.Done()

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

type CronTaskManager struct{}

func (m *CronTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup, shutdownManager shutdown.ShutdownManager) {
	defer wg.Done()

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
