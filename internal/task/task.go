package task

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskType int

const (
	ColorReset  = "\033[0m"
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func printColored(message string, colorCode string) {
	fmt.Printf("%s%s%s\n", colorCode, message, ColorReset)
}

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
	Start(ctx context.Context, task Task, wg *sync.WaitGroup)
}
