package task

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type IntervalTaskManager struct{}

func (m *IntervalTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup) {
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
