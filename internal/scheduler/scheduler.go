package scheduler

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
	ID       int           `json:"id"`
	Type     TaskType      `json:"type"`
	Name     string        `json:"name"`
	Interval time.Duration `json:"interval"`
	CronExpr string        `json:"cronExpr"`
	Disabled bool          `json:"disabled"`
	stop     chan struct{} `json:"-"`
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

func (s *Scheduler) Create(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Tasks[t.ID] = t
	go s.InvokeTask(context.Background(), t)
}

func (s *Scheduler) Read(id int) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.Tasks[id]
	return task, ok
}

func (s *Scheduler) Update(id int, t *Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.Tasks[id]; ok {
		// don't update id
		s.Tasks[id] = t

		s.stopRoutine(task)
		go s.InvokeTask(context.Background(), t)
		return true
	}

	return false
}

func (s *Scheduler) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.Tasks[id]; ok {
		s.stopRoutine(task)
		delete(s.Tasks, id)
		return true
	}

	return false
}

func (s *Scheduler) stopRoutine(task *Task) {
	// don't mutex lock here, otherwise deadlock will occur
	if task != nil {
		close(task.stop)
		// select {
		// case <-task.stop:
		// 	return

		// default:
		// 	close(task.stop)
		// }
	}
}

func (s *Scheduler) InvokeTask(ctx context.Context, task *Task) {
	task.stop = make(chan struct{})

	s.Wg.Add(1)
	defer s.Wg.Done()

	switch task.Type {
	case IntervalTask:
		s.runIntervalTask(ctx, task)

	case CronTask:
		s.runCronTask(ctx, task)
	}

	util.PrintColored(fmt.Sprintf("Task %d: Routine started...\n", task.ID), util.ColorBlue)
}

func (s *Scheduler) runIntervalTask(ctx context.Context, task *Task) {
	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-task.stop:
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorRed)
			return

		case <-ticker.C:
			if !task.Disabled {
				s.processTask(task)
			}

		case <-ctx.Done():
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorRed)
			return
		}
	}
}

func (s *Scheduler) runCronTask(ctx context.Context, task *Task) {
	c := cron.New()
	defer c.Stop()

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
	case <-task.stop:
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorRed)
		return

	case <-ctx.Done():
		util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorYellow)
		return
	}
}

func (s *Scheduler) processTask(task *Task) {
	util.PrintColored(fmt.Sprintf("Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorGreen)
	// Add your task-specific logic here
	// If an error occurs during the task execution, handle it accordingly
	// For example, errCh <- fmt.Errorf("Task %d failed", task.ID)
}
