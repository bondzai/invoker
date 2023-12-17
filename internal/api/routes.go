package api

import "net/http"

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
