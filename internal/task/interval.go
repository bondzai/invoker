package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/gracefulshutdown"
)

type IntervalTaskManager struct{}

func (m *IntervalTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup, shutdownManager gracefulshutdown.Manager) {
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
