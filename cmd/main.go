package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bondzai/goez/toolbox"
	"github.com/bondzai/invoker/internal/api"
	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/bondzai/invoker/internal/util"
)

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		fmt.Println("Initializing...")
		fmt.Println("Number of CPU:", runtime.GOMAXPROCS(runtime.NumCPU()))
		fmt.Println("Number of Goroutines:", runtime.NumGoroutine())
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	si := scheduler.NewScheduler()

	si.Tasks[1] = &scheduler.Task{
		ID:       1,
		Type:     scheduler.IntervalTask,
		Name:     "Task1",
		Interval: time.Duration(4) * time.Second,
		CronExpr: "* * * * *",
		Disabled: false,
	}

	si.Tasks[2] = &scheduler.Task{
		ID:       2,
		Type:     scheduler.IntervalTask,
		Name:     "Task2",
		Interval: time.Duration(4) * time.Second,
		CronExpr: "* * * * *",
		Disabled: false,
	}
	fmt.Println("*** Tasks ***")
	toolbox.PPrint(si.Tasks)

	go util.HandleGracefulShutdown(cancel, &si.Wg)

	server := api.NewHttpServer(si)
	go server.Start(ctx)

	for _, t := range si.Tasks {
		go si.InvokeTask(ctx, t)
	}

	si.Wg.Wait()
	select {}
}
