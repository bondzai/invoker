package main

import (
	"context"
	"sync"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/bondzai/invoker/internal/util"
)

func main() {

	var schedulerInstance *scheduler.Scheduler
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	go util.HandleGracefulShutdown(cancel, &wg)

	schedulerInstance = scheduler.NewScheduler()
	schedulerInstance.GenerateTasks(3)

	server := api.NewHttpServer(schedulerInstance)
	go server.Start(ctx)

	taskManagers := *scheduler.NewTaskManagers()

	for _, t := range schedulerInstance.Tasks {
		wg.Add(1)
		go func(task scheduler.Task) {
			defer wg.Done()
			taskManagers[task.Type].Start(ctx, task, &wg, nil)
		}(*t)
	}

	wg.Wait()

	select {}
}
