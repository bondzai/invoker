package main

import (
	"fmt"
	"time"

	"github.com/bondzai/goez/toolbox"
	"github.com/bondzai/invoker/internal/scheduler"
)

func mockTasks() map[int]*scheduler.Task {
	tasks := make(map[int]*scheduler.Task)

	tasks[1] = &scheduler.Task{
		ID:       1,
		Type:     scheduler.IntervalTask,
		Name:     "Task1",
		Interval: time.Duration(4) * time.Second,
		CronExpr: "* * * * *",
		Disabled: false,
	}

	tasks[2] = &scheduler.Task{
		ID:       2,
		Type:     scheduler.IntervalTask,
		Name:     "Task2",
		Interval: time.Duration(4) * time.Second,
		CronExpr: "* * * * *",
		Disabled: false,
	}

	fmt.Println("Mock tasks:")
	toolbox.PPrint(tasks)

	return tasks
}
