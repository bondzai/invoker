package mock

import (
	"time"

	"github.com/bondzai/invoker/internal/task"
)

func GenerateTasks(numTasks int) []task.Task {
	tasks := make([]task.Task, numTasks)

	for i := 0; i < numTasks; i++ {
		if i%2 == 0 {
			tasks[i] = task.Task{
				ID:       i + 1,
				Type:     task.IntervalTask,
				Interval: 5 * time.Second,
			}
		} else {
			tasks[i] = task.Task{
				ID:       i + 1,
				Type:     task.CronTask,
				CronExpr: "*/10 * * * *",
			}
		}
	}

	return tasks
}
