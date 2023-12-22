package scheduler

import (
	"fmt"
	"time"
)

func MockTasks() map[int]*Task {
	mode := "dynamic"
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

	maxTasks := 100000

	if mode == "dynamic" {
		for i := 1; i <= maxTasks; i++ {
			var taskType TaskType
			var projectID int
			var organization string

			if i%2 == 0 {
				projectID = 1
				taskType = IntervalTask
			} else {
				projectID = 2
				taskType = CronTask
			}

			if i < maxTasks/2 {
				organization = "org1"
			} else {
				organization = "org2"
			}

			task := &Task{
				ID:           i,
				Organization: organization,
				ProjectID:    projectID,
				Type:         taskType,
				Name:         fmt.Sprintf("Task%d", i),
				Interval:     time.Duration(i) * time.Second,
				CronExpr:     []string{"* * * * *"},
				Disabled:     false,
			}

			tasks[i] = task
		}
	}

	return tasks
}
