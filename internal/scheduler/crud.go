package scheduler

import "context"

func (s *Scheduler) Create(newTask *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Tasks[newTask.ID] = newTask
	go s.InvokeTask(context.Background(), newTask)
}

func (s *Scheduler) Read(id int) (*Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.Tasks[id]
	return task, ok
}

func (s *Scheduler) Update(id int, newTask *Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.Tasks[id]; ok {
		s.Tasks[id].Name = newTask.Name
		s.Tasks[id].Type = newTask.Type
		s.Tasks[id].Interval = newTask.Interval
		s.Tasks[id].CronExpr = newTask.CronExpr
		s.Tasks[id].Disabled = newTask.Disabled

		s.stopRoutine(task)
		go s.InvokeTask(context.Background(), newTask)
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
