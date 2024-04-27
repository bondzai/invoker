package main

import (
	"context"
	"flag"
	"log"

	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/bondzai/invoker/internal/util"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	si := scheduler.NewScheduler()

	flag.Parse()

	mockFlag := flag.Bool("mock", false, "Create dummy tasks for scheduler")
	if *mockFlag {
		si.Tasks = scheduler.MockTasks()
	}

	go util.HandleGracefulShutdown(cancel, &si.Wg)

	server := api.NewHttpServer(si)
	go server.Start(ctx)

	for _, t := range si.Tasks {
		go si.StartTask(ctx, t)
	}

	si.Wg.Wait()
	select {}
}
