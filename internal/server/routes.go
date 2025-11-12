package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) registerRoutes() {
	// Middlewares
	s.router.Use(middleware.Logger)     // لاگ کردن درخواست‌ها
	s.router.Use(middleware.Recoverer) // جلوگیری از کرش

	// مسیر Health Check
	s.router.Get("/health", s.handleHealthCheck())
	
	// گروه مسیرهای احراز هویت
	s.router.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", s.handleRegister())
		// r.Post("/login", s.handleLogin()) // (بعداً اضافه می‌شود)
	})
}