package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/bondzai/invoker/internal/mock"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port       int
	updatedMux sync.Mutex
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

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Invoker is running...")
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getTasks(w, r)
	case http.MethodPut:
		s.updateTasks(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getTasks(w http.ResponseWriter, r *http.Request) {
	s.updatedMux.Lock()
	defer s.updatedMux.Unlock()

	json.NewEncoder(w).Encode(mock.Tasks)
}

func (s *Server) updateTasks(w http.ResponseWriter, r *http.Request) {
	s.updatedMux.Lock()
	defer s.updatedMux.Unlock()

	// Parse the request body to get the updated tasks
	err := json.NewDecoder(r.Body).Decode(&mock.Tasks)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update each task
	for _, updatedTask := range *mock.Tasks {
		_ = mock.UpdateTaskWithPointer(&updatedTask)
	}

	// Respond with the updated tasks
	json.NewEncoder(w).Encode(mock.Tasks)
}
