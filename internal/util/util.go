package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
)

const (
	ColorReset  = "\033[0m"
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	TimeFormat  = "2006-01-02 15:04:05"
)

func PrintColored(message string, colorCode string) {
	fmt.Printf("%s%s%s\n", colorCode, message, ColorReset)
}

func HandleGracefulShutdown(cancel context.CancelFunc, wg *sync.WaitGroup) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	<-sigCh

	cancel()

	wg.Wait()

	log.Println("Shutdown complete.")
	os.Exit(0)
}
