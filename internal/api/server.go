package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port int
}

func NewHttpServer() *Server {
	return &Server{
		Port: 8080,
	}
}

var schedulerInstance *scheduler.Scheduler

func (s *Server) Start(ctx context.Context) error {
	schedulerInstance = scheduler.NewScheduler()

	http.HandleFunc("/tasks", s.handleTasks)

	http.Handle("/metrics", promhttp.Handler())

	serverAddr := fmt.Sprintf(":%d", s.Port)
	server := &http.Server{Addr: serverAddr}

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down the server...")
		server.Shutdown(context.Background())
	}()

	fmt.Printf("Server is listening on port %d...\n", s.Port)
	return server.ListenAndServe()
}
