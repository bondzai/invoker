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
		stop:     make(chan bool),
	}
}

func (s *Scheduler) GenerateTasks(count int) {
	for i := 1; i <= count; i++ {
		var taskType TaskType
		if i%2 == 0 {
			taskType = IntervalTask
		} else {
			taskType = CronTask
		}

		task := generateOneTask(
			i,
			taskType,
			fmt.Sprintf("Task%d", i),
			time.Duration(i)*time.Second,
			"* * * * *",
			false)

		s.Create(task)
	}
}
