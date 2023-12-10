package signalhandler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type SignalHandler struct {
	signalCh chan os.Signal
}

func NewSignalHandler() *SignalHandler {
	return &SignalHandler{
		signalCh: make(chan os.Signal, 1),
	}
}

func (sh *SignalHandler) Start() {
	signal.Notify(sh.signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sh.signalCh
		fmt.Println("Received interrupt signal. Stopping the program...")
		sh.Shutdown()
	}()
}

func (sh *SignalHandler) Shutdown() {
	close(sh.signalCh)
}
