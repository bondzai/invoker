package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bondzai/invoker/internal/scheduler"
)

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id != "" {
			s.GetTaskHandler(w, r)
		} else {
			s.GetTasksHandler(w, r)
		}

	case http.MethodPost:
		s.CreateTaskHandler(w, r)

	case http.MethodPut:
		s.UpdateTaskHandler(w, r)

	case http.MethodDelete:
		s.DeleteTaskHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, ok := schedulerInstance.Read(newTask.ID); ok {
		http.Error(w, "Task with the same ID already exists", http.StatusConflict)
		return
	}

	schedulerInstance.Create(&newTask)

	successMessage := map[string]string{"message": "Task created successfully"}
	response, err := json.Marshal(successMessage)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (s *Server) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedulerInstance.Tasks)
}

func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, ok := schedulerInstance.Read(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	response, _ := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updatedTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, ok := schedulerInstance.Read(id); !ok {
		http.NotFound(w, r)
		return
	}

	schedulerInstance.Update(id, &updatedTask)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if !schedulerInstance.Delete(id) {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
