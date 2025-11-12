package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// handleHealthCheck (نسخه بازآرایی شده)
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverStatus := "ok"
		dbStatus := "ok"
		dbError := ""

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// استفاده از pool تزریق شده
		if err := s.db.Ping(ctx); err != nil {
			dbStatus = "error"
			dbError = err.Error()
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		
		response := map[string]string{
			"status":   serverStatus,
			"database": dbStatus,
			"db_error": dbError,
			"service":  "team-app-backend",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}


// --- بخش جدید برای ثبت نام ---

type RegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

var validate = validator.New()

// handleRegister تابع ما برای ثبت نام
func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegistrationRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(req); err != nil {
			http.Error(w, "Invalid data: "+err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// صدا زدن لایه Store
		newUser, err := s.userStore.CreateUser(req.Email, string(hashedPassword))
		if err != nil {
			// TODO: ارور ایمیل تکراری را بهتر مدیریت کنیم
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}