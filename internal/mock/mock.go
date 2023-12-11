package mock

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bondzai/invoker/internal/task"
)

var (
	initialized  bool
	initMutex    sync.Mutex
	initErrorMsg = "Mock package already initialized"
	tasks        []task.Task
	numTasks     = 100000
)

func init() {
	initMutex.Lock()
	defer initMutex.Unlock()

	if initialized {
		panic(initErrorMsg)
	}

	tasks = generateTasks(numTasks)
	initialized = true
	log.Println("Mock package initialized.")
}

func generateTasks(numTasks int) []task.Task {
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
				CronExpr: "* * * * *",
			}
		}
	}

	return tasks
}

func GetTasks() []task.Task {
	return tasks
}
