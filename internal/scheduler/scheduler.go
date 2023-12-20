package scheduler

import (
	"context"
	"fmt"
	"log"
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
	IntervalTask TaskType = iota + 1
	CronTask
)

type Task struct {
	ID       int           `json:"id"`
	Type     TaskType      `json:"type"`
	Name     string        `json:"name"`
	Interval time.Duration `json:"interval"`
	CronExpr []string      `json:"cronExpr"`
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

func (s *Scheduler) stopTask(task *Task) {
	if task == nil {
		return
	}

	// don't mutex lock here, otherwise deadlock will occur
	select {
	case <-task.isAlive:
		// Channel is already closed
	default:
		close(task.isAlive)
	}
}

func (s *Scheduler) StartTask(ctx context.Context, task *Task) {
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
		s.stopTask(task)
	}()

	switch task.Type {
	case IntervalTask:
		s.runIntervalTask(ctx, task)

	case CronTask:
		s.runCronTask(ctx, task)
	}
}

func (s *Scheduler) runIntervalTask(ctx context.Context, task *Task) error {
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
			return nil

		case <-ctx.Done():
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorYellow)
			return nil
		}
	}
}

func (s *Scheduler) runCronTask(ctx context.Context, task *Task) error {
	c := cron.New()
	defer func() {
		c.Stop()
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopped\n", task.ID), util.ColorPurple)
	}()

	for _, expr := range task.CronExpr {
		localExpr := expr
		if _, err := c.AddFunc(localExpr, func() {
			if !task.Disabled {
				log.Println("Cron Task", task.ID, "Triggered with syntax: ", localExpr)
				s.processTask(task)
			}
		}); err != nil {
			log.Println("Cron Task", task.ID, "Error: ", err)
			return util.ErrInvalidTaskCronExpr
		}
	}

	c.Start()

	select {
	case <-task.isAlive:
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorRed)
		return nil

	case <-ctx.Done():
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorYellow)
		return nil
	}
}

func (s *Scheduler) processTask(task *Task) error {
	if task.Type == IntervalTask {
		util.PrintColored(fmt.Sprintf("Interval Task %d: Triggered at %v\n", task.ID, time.Now().Format(util.TimeFormat)), util.ColorGreen)
	}

	if task.Type == CronTask {
		util.PrintColored(fmt.Sprintf("Cron Task %d: Triggered at %v\n", task.ID, time.Now().Format(util.TimeFormat)), util.ColorCyan)
	}

	// Add your task-specific logic here
	// If an error occurs during the task execution, handle it accordingly
	// For example, errCh <- fmt.Errorf("Task %d failed", task.ID)
	return nil
}
