package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/util"
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
	Disabled bool // Flag to disable/enable the task
}

type TaskManager interface {
	Start(ctx context.Context, task Task, wg *sync.WaitGroup)
}

type IntervalTaskManager struct{}

func (m *IntervalTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			util.PrintColored(fmt.Sprintf("Interval Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorGreen)
			// Add your interval task-specific logic here

		case <-ctx.Done():
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorRed)
			return
		}
	}
}

type CronTaskManager struct{}

func (m *CronTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup) {
	defer wg.Done()

	c := cron.New()
	defer c.Stop()

	_, err := c.AddFunc(task.CronExpr, func() {
		util.PrintColored(fmt.Sprintf("Cron Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorPurple)
		// Add your cron task-specific logic here
	})
	if err != nil {
		util.PrintColored(fmt.Sprintf("Cron Task %d: Error adding cron expression %v\n", task.ID, err), util.ColorRed)
		return
	}

	c.Start()

	<-ctx.Done()
	util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorYellow)
}
