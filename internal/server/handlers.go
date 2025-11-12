package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverStatus := "ok"
		dbStatus := "ok"
		dbError := ""

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

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



type RegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var validate = validator.New()

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

		newUser, err := s.userStore.CreateUser(req.Email, string(hashedPassword))
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if err := validate.Struct(req); err != nil {
			http.Error(w, "Invalid data: "+err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.userStore.GetUserByEmail(req.Email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{
			"sub": user.ID,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(s.jwtSecret))
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
	}
}

func (s *Server) handleGetUserMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(userContextKey).(int)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user, err := s.userStore.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}