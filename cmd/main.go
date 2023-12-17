package main

import (
	"context"
	"sync"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/util"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	go util.HandleGracefulShutdown(cancel, &wg)

	server := api.NewHttpServer()
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
