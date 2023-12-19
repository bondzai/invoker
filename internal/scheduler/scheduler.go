package scheduler

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	ID       int           `json:"id"`
	Type     TaskType      `json:"type"`
	Name     string        `json:"name"`
	Interval time.Duration `json:"interval"`
	CronExpr string        `json:"cronExpr"`
	Disabled bool          `json:"disabled"`
	isAlive  chan struct{} `json:"-"`
}

type Scheduler struct {
	mu    sync.RWMutex
	Wg    sync.WaitGroup
	Tasks map[int]*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Tasks: make(map[int]*Task),
	}
}

func (s *Scheduler) InvokeTask(ctx context.Context, task *Task) {
	task.isAlive = make(chan struct{})

	s.Wg.Add(1)
	defer s.Wg.Done()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		select {
		case <-stop:
		case <-ctx.Done():
		}
		s.stopRoutine(task)
	}()

	switch task.Type {
	case IntervalTask:
		s.runIntervalTask(ctx, task)

	case CronTask:
		s.runCronTask(ctx, task)
	}
}

func (s *Scheduler) runIntervalTask(ctx context.Context, task *Task) {
	ticker := time.NewTicker(task.Interval)
	defer func() {
		ticker.Stop()
		util.PrintColored(fmt.Sprintf("Interval Task %d: Stopped\n", task.ID), util.ColorPurple)
	}()

	for {
		select {
		case <-ticker.C:
			if !task.Disabled {
				s.processTask(task)
			}

		case <-task.isAlive:
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorRed)
			return

		case <-ctx.Done():
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorYellow)
			return
		}
	}
}

func (s *Scheduler) runCronTask(ctx context.Context, task *Task) {
	c := cron.New()
	defer func() {
		c.Stop()
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopped\n", task.ID), util.ColorPurple)
	}()

	_, err := c.AddFunc(task.CronExpr, func() {
		if !task.Disabled {
			s.processTask(task)
		}
	})
	if err != nil {
		util.PrintColored(fmt.Sprintf("Cron Task %d: Error adding cron expression %v\n", task.ID, err), util.ColorRed)
		return
	}

	c.Start()

	select {
	case <-task.isAlive:
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorRed)
		return

	case <-ctx.Done():
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorYellow)
		return
	}
}

func (s *Scheduler) processTask(task *Task) {
	if task.Type == IntervalTask {
		util.PrintColored(fmt.Sprintf("Interval Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorGreen)
	} else {
		util.PrintColored(fmt.Sprintf("Cron Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorCyan)
	}

	// Add your task-specific logic here
	// If an error occurs during the task execution, handle it accordingly
	// For example, errCh <- fmt.Errorf("Task %d failed", task.ID)
}

func (s *Scheduler) stopRoutine(task *Task) {
	// don't mutex lock here, otherwise deadlock will occur
	if task != nil {
		select {
		case <-task.isAlive:
			// Channel is already closed
		default:
			close(task.isAlive)
		}
	}
}
