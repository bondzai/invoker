package scheduler

import "context"

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
		s.Tasks[id].Name = t.Name
		s.Tasks[id].Type = t.Type
		s.Tasks[id].Interval = t.Interval
		s.Tasks[id].CronExpr = t.CronExpr
		s.Tasks[id].Disabled = t.Disabled

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
