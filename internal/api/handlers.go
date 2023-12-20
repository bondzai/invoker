package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bondzai/invoker/internal/scheduler"
)

// Common error messages
const (
	InvalidRequestPayload = "Invalid request payload"
	InvalidTaskID         = "Invalid task ID"
	TaskIDNotExists       = "Task ID not exists"
	TaskConflict          = "Task with the same ID already exists"
)

// Common success messages
const (
	TaskCreatedSuccessfully = "Task created successfully"
	TaskUpdatedSuccessfully = "Task updated successfully"
	TaskDeletedSuccessfully = "Task deleted successfully"
)

func (s *Server) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.Scheduler.Tasks)
}

func (s *Server) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.response(w, http.StatusBadRequest, InvalidRequestPayload)
		return
	}

	task, ok := s.Scheduler.Read(id)
	if !ok {
		s.response(w, http.StatusNotFound, TaskIDNotExists)
		return
	}

	response, _ := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		s.response(w, http.StatusBadRequest, InvalidRequestPayload)
		return
	}

	if _, ok := s.Scheduler.Read(newTask.ID); ok {
		s.response(w, http.StatusConflict, TaskConflict)
		return
	}

	s.Scheduler.Create(&newTask)
	s.response(w, http.StatusCreated, TaskCreatedSuccessfully)
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.response(w, http.StatusBadRequest, InvalidTaskID)
		return
	}

	var updatedTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		s.response(w, http.StatusBadRequest, InvalidRequestPayload)
		return
	}

	if _, ok := s.Scheduler.Read(id); !ok {
		s.response(w, http.StatusNotFound, TaskIDNotExists)
		return
	}

	s.Scheduler.Update(id, &updatedTask)
	s.response(w, http.StatusOK, TaskUpdatedSuccessfully)
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.response(w, http.StatusBadRequest, InvalidTaskID)
		return
	}

	if !s.Scheduler.Delete(id) {
		s.response(w, http.StatusNotFound, TaskIDNotExists)
		return
	}

	s.response(w, http.StatusOK, TaskDeletedSuccessfully)
}
