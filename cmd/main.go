package main

import (
	"context"
	"sync"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/bondzai/invoker/internal/util"
)

var schedulerInstance *scheduler.Scheduler
var wg sync.WaitGroup

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go util.HandleGracefulShutdown(cancel, &wg)

	schedulerInstance = scheduler.NewScheduler()
	scheduler.GenerateTasks(schedulerInstance, 5)

	server := api.NewHttpServer(schedulerInstance)
	go server.Start(ctx)

	// taskManagers := *task.NewTaskManagers()
	// taskFromDB := mock.Tasks

	// for _, t := range *taskFromDB {
	// 	wg.Add(1)
	// 	go func(task task.Task) {
	// 		defer wg.Done()
	// 		taskManagers[task.Type].Start(ctx, task, &wg, nil)
	// 	}(t)
	// }

	wg.Wait()

	select {}
}
