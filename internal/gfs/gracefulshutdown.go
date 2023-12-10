package gracefulshutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// Manager defines a generic shutdown manager.
type Manager struct {
	signalCh  chan os.Signal
	cancelCtx context.Context
	cancel    context.CancelFunc
	wgCounter int32
	mutex     sync.Mutex
}

// NewManager creates a new graceful shutdown manager.
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		signalCh:  make(chan os.Signal, 1),
		cancelCtx: ctx,
		cancel:    cancel,
		wgCounter: 0,
		mutex:     sync.Mutex{},
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
	m.mutex.Lock()
	atomic.AddInt32(&m.wgCounter, 1)
	m.mutex.Unlock()

	go func() {
		<-m.cancelCtx.Done()

		m.mutex.Lock()
		defer m.mutex.Unlock()

		// Ensure decrementing the counter only if it's greater than 0
		if m.wgCounter > 0 {
			atomic.AddInt32(&m.wgCounter, -1)
			fmt.Printf("Decremented WaitGroup counter: %v\n", m.wgCounter)
			wg.Done()
		} else {
			fmt.Println("WaitGroup counter is already zero, skipping decrement.")
		}
	}()

	return &wg
}

// WaitGroupCounter returns the current counter value of the associated WaitGroup.
func (m *Manager) WaitGroupCounter() int32 {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return atomic.LoadInt32(&m.wgCounter)
}
