package scheduler

import (
	"fmt"
	"sync"
	"time"
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

// Create adds a new item to the storage.
func (s *Scheduler) Create(item *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Tasks[item.ID] = item
}

// Read retrieves an item from the storage by ID.
func (s *Scheduler) Read(id int) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.Tasks[id]
	return task, ok
}

// Update updates an existing item in the storage.
func (s *Scheduler) Update(id int, newItem *Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Tasks[id]; ok {
		s.Tasks[id] = newItem
		return true
	}
	return false
}

// Delete removes an item from the storage by ID.
func (s *Scheduler) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("start delete...")

	fmt.Println("Tasks before delete:", s.Tasks)

	if _, ok := s.Tasks[id]; ok {
		delete(s.Tasks, id)
		fmt.Println("delete success...")
		fmt.Println("Tasks after delete:", s.Tasks)
		return true
	}

	fmt.Println("delete failed...")
	return false
}
