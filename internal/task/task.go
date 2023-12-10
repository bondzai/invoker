package task

import (
	"context"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/shutdown"
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
