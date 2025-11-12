package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
    // مطمئن شوید نام ماژول شما در go.mod صحیح است
	"team-app-backend/internal/user" 
)

// Server struct اصلی است
type Server struct {
	router *chi.Mux
	db     *pgxpool.Pool

	userStore *user.Store
	jwtSecret	string
}

func NewServer(dbPool *pgxpool.Pool, jwtSecret string) *Server {
	
	userStore := user.NewStore(dbPool)

	s := &Server{
		router: chi.NewRouter(),
		db:     dbPool,
		userStore: userStore,
		jwtSecret: jwtSecret,
	}

	s.registerRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}