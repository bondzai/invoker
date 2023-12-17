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
	stop     chan bool     `json:"-"`
}

type Scheduler struct {
	mu    sync.RWMutex
	Tasks map[int]*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Tasks: make(map[int]*Task),
	}
}

func (s *Scheduler) Create(item *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Tasks[item.ID] = item
}

func (s *Scheduler) Read(id int) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.Tasks[id]
	return task, ok
}

func (s *Scheduler) Update(id int, newItem *Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Tasks[id]; ok {
		s.Tasks[id] = newItem
		return true
	}
	return false
}

func (s *Scheduler) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Tasks[id]; ok {
		delete(s.Tasks, id)
		return true
	}

	return false
}

type TaskManager interface {
	Start(ctx context.Context, task Task, wg *sync.WaitGroup, taskCh chan<- Task)
}

func NewTaskManagers() *map[TaskType]TaskManager {
	return &map[TaskType]TaskManager{
		IntervalTask: &IntervalTaskManager{},
		CronTask:     &CronTaskManager{},
	}
}

type IntervalTaskManager struct{}

func (m *IntervalTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup, taskCh chan<- Task) {
	defer wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !task.Disabled {
				util.PrintColored(fmt.Sprintf("Interval Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorGreen)
				// Add your interval task-specific logic here
				// If an error occurs during the task execution, send it to the error channel
				// For example, errCh <- fmt.Errorf("Interval task %d failed", task.ID)
			}

		case <-ctx.Done():
			util.PrintColored(fmt.Sprintf("Interval Task %d: Stopping...\n", task.ID), util.ColorRed)
			return
		}
	}
}

type CronTaskManager struct{}

func (m *CronTaskManager) Start(ctx context.Context, task Task, wg *sync.WaitGroup, taskCh chan<- Task) {
	defer wg.Done()

	c := cron.New()
	defer c.Stop()

	_, err := c.AddFunc(task.CronExpr, func() {
		if !task.Disabled {
			util.PrintColored(fmt.Sprintf("Cron Task %d: Triggered at %v\n", task.ID, time.Now().Format(time.RFC3339)), util.ColorPurple)
			// Add your cron task-specific logic here
			// If an error occurs during the task execution, send it to the error channel
			// For example, errCh <- fmt.Errorf("Cron task %d failed", task.ID)
		}
	})
	if err != nil {
		util.PrintColored(fmt.Sprintf("Cron Task %d: Error adding cron expression %v\n", task.ID, err), util.ColorRed)
		// Send the error to the channel
		// errCh <- err
		return
	}

	c.Start()

	<-ctx.Done()
	util.PrintColored(fmt.Sprintf("Cron Task %d: Stopping...\n", task.ID), util.ColorYellow)
}
