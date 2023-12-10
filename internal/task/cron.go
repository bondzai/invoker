package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type CronTaskManager struct{}

func (m *CronTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup) {
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
