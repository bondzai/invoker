package scheduler

import (
	"fmt"
	"time"
)

func generateOneTask(id int, taskType TaskType, name string, interval time.Duration, cronExpr string, disabled bool) *Task {
	return &Task{
		ID:       id,
		Type:     taskType,
		Name:     name,
		Interval: interval,
		CronExpr: cronExpr,
		Disabled: disabled,
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) GenerateTasks(count int) {
	for i := 1; i <= count; i++ {
		task := generateOneTask(
			i,
			IntervalTask,
			fmt.Sprintf("Task%d", i),
			time.Duration(4)*time.Second,
			"* * * * *",
			false)

		s.Create(task)
	}
}
