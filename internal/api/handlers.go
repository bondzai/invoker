package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bondzai/invoker/internal/scheduler"
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
		s.response(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	task, ok := s.Scheduler.Read(id)
	if !ok {
		s.response(w, http.StatusNotFound, "Task ID not exists")
		return
	}

	response, _ := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *Server) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		s.response(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if _, ok := s.Scheduler.Read(newTask.ID); ok {
		s.response(w, http.StatusConflict, "Task with the same ID already exists")
		return
	}

	s.Scheduler.Create(&newTask)
	s.response(w, http.StatusCreated, "Task created successfully")
}

func (s *Server) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.response(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var updatedTask scheduler.Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		s.response(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if _, ok := s.Scheduler.Read(id); !ok {
		s.response(w, http.StatusNotFound, "Task ID not exists")
		return
	}

	s.Scheduler.Update(id, &updatedTask)
	s.response(w, http.StatusOK, "Task updated successfully")
}

func (s *Server) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.response(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if !s.Scheduler.Delete(id) {
		s.response(w, http.StatusNotFound, "Task ID not exists")
		return
	}

	s.response(w, http.StatusOK, "Task deleted successfully")
}
