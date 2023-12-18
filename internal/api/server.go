package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bondzai/invoker/internal/scheduler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Port      int
	Scheduler *scheduler.Scheduler
}

func NewHttpServer(scheduler *scheduler.Scheduler) *Server {
	return &Server{
		Port:      8080,
		Scheduler: scheduler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	http.HandleFunc("/tasks", s.handleTasks)

	http.Handle("/metrics", promhttp.Handler())

	serverAddr := fmt.Sprintf(":%d", s.Port)
	server := &http.Server{Addr: serverAddr}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down the server...")
		server.Shutdown(context.Background())
	}()

	log.Printf("Server is listening on port %d...\n", s.Port)
	return server.ListenAndServe()
}

func (s *Server) sendResponseMessage(w http.ResponseWriter, statusCode int, message string) {
	responseData := map[string]string{"message": message}
	response, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
