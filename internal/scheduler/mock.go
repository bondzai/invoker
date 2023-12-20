package scheduler

import (
	"fmt"
	"time"
)

func MockTasks() map[int]*Task {
	mode := "static"
	tasks := make(map[int]*Task)

	if mode == "static" {
		tasks[1] = &Task{
			ID:       1,
			Type:     IntervalTask,
			Name:     "Task1",
			Interval: time.Duration(5) * time.Second,
			CronExpr: []string{"* * * * *"},
			Disabled: false,
		}

		tasks[2] = &Task{
			ID:       2,
			Type:     CronTask,
			Name:     "Task2",
			Interval: time.Duration(4) * time.Second,
			CronExpr: []string{"*/2 * * * *", "* * * * *"},
			Disabled: false,
		}
	}

	if mode == "dynamic" {
		for i := 1; i <= 100000; i++ {
			var taskType TaskType
			if i%2 == 0 {
				taskType = IntervalTask
			} else {
				taskType = CronTask
			}

			task := &Task{
				ID:       i,
				Type:     taskType,
				Name:     fmt.Sprintf("Task%d", i),
				Interval: time.Duration(i) * time.Second,
				CronExpr: []string{"* * * * *"},
				Disabled: false,
			}
			tasks[i] = task
		}
	}

	return tasks
}
