package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey
type contextKey string
const userContextKey = contextKey("userID")

// AuthMiddleware
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ۱. read Header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// ۲."Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// ۳. Token validation
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(s.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := (*claims)["sub"].(float64)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat) // conver to int

		ctx := context.WithValue(r.Context(), userContextKey, userID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}