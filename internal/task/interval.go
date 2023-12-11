package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/util"
)

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
