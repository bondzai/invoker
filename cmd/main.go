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

	s := scheduler.NewScheduler()

	flag.Parse()

	mockFlag := flag.Bool("mock", false, "Create dummy tasks for scheduler")
	if *mockFlag {
		s.Tasks = scheduler.MockTasks()
	}

	go util.HandleGracefulShutdown(cancel, &s.Wg)

	server := api.NewHttpServer(s)
	go server.Start(ctx)

	for _, t := range s.Tasks {
		go s.StartTask(ctx, t)
	}

	s.Wg.Wait()
	select {}
}
