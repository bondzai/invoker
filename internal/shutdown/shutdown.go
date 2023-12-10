// shutdown.go
package shutdown

import (
	"context"
)

type ShutdownManager interface {
	Shutdown()
}

type GracefulShutdownManager struct {
	cancelFunc context.CancelFunc
}

func NewGracefulShutdownManager() *GracefulShutdownManager {
	_, cancel := context.WithCancel(context.Background())
	return &GracefulShutdownManager{
		cancelFunc: cancel,
	}
}

func (m *GracefulShutdownManager) Shutdown() {
	m.cancelFunc()
}
