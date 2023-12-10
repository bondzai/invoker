package gracefulshutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// Manager defines a generic shutdown manager.
type Manager struct {
	signalCh  chan os.Signal
	cancelCtx context.Context
	cancel    context.CancelFunc
}

// NewManager creates a new graceful shutdown manager.
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		signalCh:  make(chan os.Signal, 1),
		cancelCtx: ctx,
		cancel:    cancel,
	}
}

// StartSignalHandling starts listening for interrupt and SIGTERM signals.
// When a signal is received, it triggers the graceful shutdown.
func (m *Manager) StartSignalHandling() {
	signal.Notify(m.signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-m.signalCh
		fmt.Println("Received interrupt signal. Stopping the program...")
		m.Shutdown()
	}()
}

// Shutdown cancels the context, triggering a graceful shutdown.
func (m *Manager) Shutdown() {
	m.cancel()
}

// Context returns the cancellation context associated with the manager.
func (m *Manager) Context() context.Context {
	return m.cancelCtx
}

// WaitGroup returns a new sync.WaitGroup that is associated with the cancellation context.
func (m *Manager) WaitGroup() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		<-m.cancelCtx.Done()
		wg.Done()
	}()
	return &wg
}
