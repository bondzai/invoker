package main

import (
	"context"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/bondzai/invoker/internal/util"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	si := scheduler.NewScheduler()
	si.GenerateTasks(3)

	go util.HandleGracefulShutdown(cancel, &si.Wg)

	server := api.NewHttpServer(si)
	go server.Start(ctx)

	go func() {
		for _, t := range si.Tasks {
			si.Wg.Add(1)
			go si.InvokeTask(ctx, t)
		}
	}()

	si.Wg.Wait()
	select {}
}
