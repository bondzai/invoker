package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/bondzai/invoker/internal/util"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	si := scheduler.NewScheduler()

	si.Tasks = mockTasks()

	go util.HandleGracefulShutdown(cancel, &si.Wg)

	server := api.NewHttpServer(si)
	go server.Start(ctx)

	for _, t := range si.Tasks {
		go si.InvokeTask(ctx, t)
	}

	si.Wg.Wait()
	select {}
}

func mockTasks() map[int]*scheduler.Task {
	mode := "static"
	tasks := make(map[int]*scheduler.Task)

	if mode == "static" {
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
			Type:     scheduler.CronTask,
			Name:     "Task2",
			Interval: time.Duration(4) * time.Second,
			CronExpr: "* * * * *",
			Disabled: false,
		}
	}

	if mode == "dynamic" {
		for i := 1; i <= 100000; i++ {
			var taskType scheduler.TaskType
			if i%2 == 0 {
				taskType = scheduler.IntervalTask
			} else {
				taskType = scheduler.CronTask
			}

			task := &scheduler.Task{
				ID:       i,
				Type:     taskType,
				Name:     fmt.Sprintf("Task%d", i),
				Interval: time.Duration(i) * time.Second,
				CronExpr: "* * * * *",
				Disabled: false,
			}
			tasks[i] = task
		}
	}

	return tasks
}
