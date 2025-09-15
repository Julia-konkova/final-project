package server

import (
	"final-project/handlers"
	"final-project/pkg/api"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router *chi.Mux
	port   string
}

func New(port string) *Server {
	return &Server{
		router: chi.NewRouter(),
		port:   port,
	}
}

func (s *Server) SetupRoutes(h *handlers.Handlers) {

	s.router.Get("/*", h.ServeWebFiles)

	s.router.Get("/api/nextdate", api.NextDayHandler)

	s.router.Post("/api/task", h.AddTaskHandler)

	s.router.Get("/api/tasks", h.TasksHandler)

	s.router.Get("/api/task", h.GetTaskHandler)

	s.router.Put("/api/task", h.UpdateTaskHandler)

	s.router.Post("/api/task/done", h.DoneTaskHandler)

	s.router.Delete("/api/task", h.DeleteTaskHandler)
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.port, s.router)
}

func (s *Server) Router() *chi.Mux {
	return s.router
}
