package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bondzai/invoker/internal/mock"
)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mock.Tasks)
}

func (s *Server) updateTasks(w http.ResponseWriter, r *http.Request) {
	s.updatedMux.Lock()
	defer s.updatedMux.Unlock()

	err := json.NewDecoder(r.Body).Decode(&mock.Tasks)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, updatedTask := range *mock.Tasks {
		_ = mock.UpdateTaskWithPointer(&updatedTask)
		fmt.Println("Updated task", updatedTask)
		fmt.Println("Sent to channel successfully")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mock.Tasks)
}
