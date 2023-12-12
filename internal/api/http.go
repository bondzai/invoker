package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bondzai/invoker/internal/mock"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port int
}

// NewServer creates a new Server instance with default values
func NewHttpServer() *Server {
	return &Server{
		Port: 8080,
	}
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	http.HandleFunc("/ping", s.pingHandler)

	http.HandleFunc("/tasks", s.getTasks)

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

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Invoker is running...")
}

func (s *Server) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks := mock.GetStaticTasks()
	json.NewEncoder(w).Encode(tasks)
}
