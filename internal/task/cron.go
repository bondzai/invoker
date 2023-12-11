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
		printColored(fmt.Sprintf("Cron Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), ColorPurple)
		// Add your cron task-specific logic here
	})
	if err != nil {
		printColored(fmt.Sprintf("Cron Task %d: Error adding cron expression %v\n", task.ID, err), ColorRed)
		return
	}

	c.Start()

	<-ctx.Done()
	printColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), ColorYellow)
}
