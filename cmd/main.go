package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a channel to signal task execution
	taskChannel := make(chan bool)

	// Start interval task in a separate goroutine
	go intervalTask(taskChannel, 5*time.Second)

	// Start cron task in a separate goroutine
	go cronTask(taskChannel, "0 0 * * *")

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Signal tasks to stop
	close(taskChannel)

	fmt.Println("Shutting down gracefully")
	time.Sleep(2 * time.Second) // Give some time for tasks to finish gracefully
}

func intervalTask(ch <-chan bool, interval time.Duration) {
	for {
		select {
		case <-ch:
			fmt.Println("Interval task stopped")
			return
		default:
			fmt.Println("Executing interval task")
			// Your interval task logic goes here

			time.Sleep(interval)
		}
	}
}

func cronTask(ch <-chan bool, cronExpr string) {
	cronSchedule, err := time.Parse("cron", cronExpr)
	if err != nil {
		fmt.Println("Error parsing cron expression:", err)
		return
	}

	for {
		select {
		case <-ch:
			fmt.Println("Cron task stopped")
			return
		default:
			now := time.Now()
			nextRun := cronSchedule
			if now.After(nextRun) {
				nextRun = cronSchedule.Add(cronSchedule.Sub(now))
			}

			sleepTime := nextRun.Sub(now)
			fmt.Printf("Waiting for %s until next cron task\n", sleepTime)
			time.Sleep(sleepTime)

			fmt.Println("Executing cron task")
			// Your cron task logic goes here
		}
	}
}
