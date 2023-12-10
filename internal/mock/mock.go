package mock

import (
	"math/rand"
	"time"

	"github.com/bondzai/invoker/internal/task"
)

func GenerateTasks(numTasks int) []task.Task {
	tasks := make([]task.Task, numTasks)

	for i := 0; i < numTasks; i++ {
		if i%2 == 0 {
			interval := time.Duration(rand.Intn(56)+5) * time.Second

			tasks[i] = task.Task{
				ID:       i + 1,
				Type:     task.IntervalTask,
				Interval: interval,
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
