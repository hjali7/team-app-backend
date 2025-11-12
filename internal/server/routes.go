package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) registerRoutes() {
	// Middlewares
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://team.looparc.ir"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	}))
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/health", s.handleHealthCheck())
	
	s.router.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", s.handleRegister())
		r.Post("/login", s.handleLogin()) 
	})
}