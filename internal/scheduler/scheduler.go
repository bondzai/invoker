package scheduler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	ID       int
	Name     string
	Interval time.Duration
	Action   func()
	stop     chan bool
}

type Scheduler struct {
	Tasks []*Task
	mu    sync.Mutex // Mutex for safe concurrent access to the Tasks
}

func (s *Scheduler) AddTask(task *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Tasks = append(s.Tasks, task)
}

func (s *Scheduler) RunScheduler() {
	var wg sync.WaitGroup

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.Tasks {
		wg.Add(1)
		go func(t *Task) {
			defer wg.Done()
			for {
				select {
				case <-time.After(t.Interval):
					fmt.Println("")
					fmt.Println(time.Now().Format("15:04:05"))
					fmt.Println(t.Interval)
					t.Action()
					fmt.Println("")
				case <-t.stop:
					return
				}
			}
		}(task)
	}

	wg.Wait()
}

func (s *Scheduler) UpdateTaskInterval(taskID int, newInterval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.Tasks {
		if task.ID == taskID {
			task.stop <- true

			task.Interval = newInterval

			go func(t *Task) {
				for {
					select {
					case <-time.After(t.Interval):
						fmt.Println("")
						fmt.Println(time.Now().Format("15:04:05"))
						fmt.Println(t.Interval)
						t.Action()
						fmt.Println("")
					case <-t.stop:
						return
					}
				}
			}(task)

			return
		}
	}
}

func updateTaskHandler(s *Scheduler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			TaskID      int    `json:"task_id"`
			NewInterval string `json:"new_interval"`
		}

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding request body: %s", err), http.StatusBadRequest)
			return
		}

		duration, err := time.ParseDuration(requestData.NewInterval)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing duration: %s", err), http.StatusBadRequest)
			return
		}

		s.UpdateTaskInterval(requestData.TaskID, duration)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Task %d interval updated to %s", requestData.TaskID, requestData.NewInterval)
	}
}

func Scheduling() {
	scheduler := &Scheduler{}

	task1 := &Task{
		ID:       1,
		Name:     "Task 1",
		Interval: 2 * time.Second,
		Action: func() {
			fmt.Println("Running Task 1")
		},
		stop: make(chan bool),
	}

	task2 := &Task{
		ID:       2,
		Name:     "Task 2",
		Interval: 5 * time.Second,
		Action: func() {
			fmt.Println("Running Task 2")
		},
		stop: make(chan bool),
	}

	scheduler.AddTask(task1)
	scheduler.AddTask(task2)

	fmt.Println("Scheduler is running...")

	go scheduler.RunScheduler()

	http.HandleFunc("/updateTask", updateTaskHandler(scheduler))

	http.ListenAndServe(":8080", nil)
}
