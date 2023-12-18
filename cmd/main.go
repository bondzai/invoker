package main

import (
	"context"
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
		Type:     scheduler.CronTask,
		Name:     "Task2",
		Interval: time.Duration(4) * time.Second,
		CronExpr: "* * * * *",
		Disabled: false,
	}

	return tasks
}
