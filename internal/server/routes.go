package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) registerRoutes() {
	// Middlewares
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/health", s.handleHealthCheck())
	
	s.router.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", s.handleRegister())
		r.Post("/login", s.handleLogin()) 
	})
}